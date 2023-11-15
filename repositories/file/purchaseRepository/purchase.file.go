package purchaseFileRepository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	model "haf.systems/purchases/models/purchase"

	"haf.systems/purchases/utils"
)

type purchaseJsonEntity struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Value       float32   `json:"value"`
}

func (p purchaseJsonEntity) ToPurchase() model.Purchase {
	return model.NewPurchase(p.Description, p.Date, p.Value)
}

func fromPurchase(id string, p model.Purchase) purchaseJsonEntity {
	return purchaseJsonEntity{
		Id:          id,
		Description: p.Description(),
		Date:        p.Date(),
		Value:       p.Value(),
	}
}

type PurchaseFileRepository struct {
	filename  string
	purchases map[string]model.Purchase
}

func loadFromFile(filename string) (*map[string]model.Purchase, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var entries []purchaseJsonEntity
	err = json.Unmarshal(fileBytes, &entries)
	if err != nil {
		return nil, err
	}

	mapped := make(map[string]model.Purchase, len(entries))
	for _, e := range entries {
		mapped[e.Id] = e.ToPurchase()
	}

	return &mapped, nil
}

func NewPurchaseFileRepository(_filename string) *PurchaseFileRepository {

	p, e := loadFromFile(_filename)

	if e != nil {
		println(e)

		pm := make(map[string]model.Purchase)
		p = &pm
	}

	return &PurchaseFileRepository{
		purchases: *p,
		filename:  _filename,
	}
}

func (repo PurchaseFileRepository) GetPurchase(id string) (*model.Purchase, *utils.HTTPError) {
	if p, exists := repo.purchases[id]; exists {
		return &p, nil
	}

	return nil, &utils.HTTPError{
		StatusCode: http.StatusNotFound,
		Err:        fmt.Errorf("not found"),
	}
}

func (repo PurchaseFileRepository) RecordPurchase(id string, p model.Purchase) error {
	// check if already recorded
	// if yes, error out
	if _, exists := repo.purchases[id]; exists {
		return utils.HTTPError{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        fmt.Errorf("entry already exists"),
		}
	}

	// lets try to append the file first, then adding to the map
	toRecord := fromPurchase(id, p)

	newBytes, err := json.Marshal(toRecord)
	if err != nil {
		return utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	file, err := os.OpenFile(repo.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(newBytes)
	if err != nil {
		return utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// finally adding to the map
	repo.purchases[id] = p

	return nil
}
