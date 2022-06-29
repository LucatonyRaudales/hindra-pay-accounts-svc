package models

import (
	"errors"
	"html"
	"strings"
	"time"

  "gorm.io/gorm"
)
type WalletAccount struct {
	ID        uint64 		`gorm:"primary_key;auto_increment" json:"id"`
	Address   string 		`gorm:"size:255;not null;unique" json:"wallet_address"`
	Coin   		string 		`gorm:"size:255;not null;" json:"coin"`
	UserID   	string 		`gorm:"size:255;not null;" json:"user_id"`
	Enabled 	bool 			`gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (wa *WalletAccount) Prepare() {
	wa.ID = 0
	wa.Address = html.EscapeString(strings.TrimSpace(wa.Address))
	wa.Coin = html.EscapeString(strings.TrimSpace(wa.Coin))
	wa.UserID = html.EscapeString(strings.TrimSpace(wa.UserID))
	wa.Enabled = true
	wa.CreatedAt = time.Now()
	wa.UpdatedAt = time.Now()
}

func (p *WalletAccount) Validate() error {

	if p.Address == "" {
		return errors.New("Required wallet address")
	}
	if p.Coin == "" {
		return errors.New("Required coin")
	}
	if p.UserID == "" {
		return errors.New("Required UserID")
	}
	return nil
}

func (wa *WalletAccount) SaveWalletAccount(db *gorm.DB) (*WalletAccount, error) {
	var err error
	err = db.Debug().Model(&WalletAccount{}).Create(&wa).Error
	if err != nil {
		return &WalletAccount{}, err
	}
	return wa, nil
}

func (p *WalletAccount) FindAllWalletAccounts(db *gorm.DB) (*[]WalletAccount, error) {
	var err error
	WalletAccounts := []WalletAccount{}
	err = db.Debug().Model(&WalletAccount{}).Limit(100).Find(&WalletAccounts).Error
	if err != nil {
		return &[]WalletAccount{}, err
	}
	return &WalletAccounts, nil
}

func (p *WalletAccount) FindWalletAccountByID(db *gorm.DB, pid uint32) (*WalletAccount, error) {
	var err error
	err = db.Debug().Model(&WalletAccount{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &WalletAccount{}, err
	}
	return p, nil
}


func (p *WalletAccount) DisableWalletAccount(db *gorm.DB, pid uint64, uid string) (int64, error) {

	db = db.Debug().Model(&WalletAccount{}).Where("id = ? and User_id = ?", pid, uid).Updates(WalletAccount{Enabled: false})

	if db.Error != nil {
		/*if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("WalletAccount not found")
		}*/
		return 0, db.Error
	}
	return db.RowsAffected, nil
}