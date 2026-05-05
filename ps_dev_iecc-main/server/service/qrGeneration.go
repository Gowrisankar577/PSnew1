package service

import (
	"math/rand"
	"ps_portal/db"
	"time"
)

func QrCodeGeneration(qrFor string, qrIds string, createdBy any) (bool, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	qrCode := r.Intn(900000) + 100000
	expireAt := time.Now().Add(time.Second * 12).Format("2006-01-02 15:04:05")
	_, err := db.DB.Exec("insert into app_qr_code (qr_code,qr_for,qr_ids,created_by,expire_at) values (?,?,?,?,?)", qrCode, qrFor, qrIds, createdBy, expireAt)
	return err == nil, qrCode
}
