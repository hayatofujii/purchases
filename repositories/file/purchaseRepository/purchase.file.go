package purchaseFileRepository

import (
	"encoding/json"
	"io"
	"os"

	model "haf.systems/purchases/models/purchase"
)

type purchaseJsonEntity struct {
	model.Purchase
	Id string `json:"id"`
}

func (p purchaseJsonEntity) ToPurchase() model.Purchase {
	return p.Purchase
}

func fromPurchase(id string, p model.Purchase) purchaseJsonEntity {
	return purchaseJsonEntity{
		Id:       id,
		Purchase: p,
	}
}

type PurchaseFileRepository struct {
	file      *os.File
	purchases map[string]model.Purchase
}

func loadFromFile(file *os.File) (map[string]model.Purchase, error) {

	fileBytes, err := io.ReadAll(file)
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

	return mapped, nil
}

func NewPurchaseFileRepository(_filename string) *PurchaseFileRepository {
	file, err := os.OpenFile(_filename, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		panic(err)
	}

	p, e := loadFromFile(file)

	if e != nil {
		println(e)

		p = make(map[string]model.Purchase)
	}

	return &PurchaseFileRepository{
		purchases: p,
		file:      file,
	}
}

func (repo PurchaseFileRepository) GetPurchase(id string) (bool, *model.Purchase) {
	if p, exists := repo.purchases[id]; exists {
		return true, &p
	}

	return false, nil
}

// Returns (true, nil) if purchase was recorded into file; (false, nil) if there's already an entry for it;
// (false, !nil) if it's not written into file due to some error.
func (repo PurchaseFileRepository) RecordPurchase(id string, p model.Purchase) (bool, error) {
	// check if already recorded
	if exists, _ := repo.GetPurchase(id); exists {
		return false, nil
	}

	err := repo.append(id, p)
	if err != nil {
		delete(repo.purchases, id)
		return false, err
	}

	repo.purchases[id] = p

	return true, nil
}

func (repo PurchaseFileRepository) append(id string, p model.Purchase) error {

	// file startup
	if len(repo.purchases) == 0 {
		_, err := repo.file.Write([]byte{'[', ']'})

		if err != nil {
			return err
		}
	}

	// set file pointer to just before final ']'
	_, err := repo.file.Seek(-1, 2)
	if err != nil {
		return err
	}

	// write a ',' if there's anything already
	if len(repo.purchases) > 0 {
		_, err := repo.file.Write([]byte{','})

		if err != nil {
			return err
		}
	}

	newBytes, err := json.Marshal(fromPurchase(id, p))
	if err != nil {
		return err
	}

	// file finisher
	newBytes = append(newBytes, []byte{']'}...)

	_, err = repo.file.Write(newBytes)
	if err != nil {
		return err
	}

	return nil
}

func (repo PurchaseFileRepository) Close() error {

	return repo.file.Close()
}
