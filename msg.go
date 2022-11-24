package main

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"time"

	cryptoRand "crypto/rand"

	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20poly1305"
)

const defaultDeviceKey = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const MaxUdpPacketLen = (1300)
const HeaderLen = 16 + 12
const MaxPayloadLen = (MaxUdpPacketLen - HeaderLen - chacha20poly1305.Overhead)


type MsgBuilder struct {
	Cipher cipher.AEAD
}

func (msg *MsgBuilder) Init() {
	cipher, err := chacha20poly1305.New([]byte(defaultDeviceKey))
	if err != nil {
		panic(err)
	}
	msg.Cipher = cipher
}

func (msg *MsgBuilder) BuildMsg(payload []byte, deviceId uuid.UUID) ([]byte, error) {
	pkt := make([]byte, 16+12+chacha20poly1305.Overhead+len(payload))
	var nonce [12]byte
	_, err := cryptoRand.Read(nonce[:])
	if err != nil {
		return nil, err
	}
	var timeStamp uint32 = uint32(time.Now().Unix())
	binary.LittleEndian.PutUint32(nonce[:4], timeStamp)
	copy(pkt[:16], deviceId[:])
	copy(pkt[16:16+12], nonce[:])
	ciphertext := msg.Cipher.Seal(nil, nonce[:], payload, nil)
	copy(pkt[16+12:], ciphertext[:])
	return pkt, nil
}

func (msg *MsgBuilder) ParseMsg(pkt []byte) ([]byte, uuid.UUID, error) {
	payloadLen := len(pkt) - HeaderLen - chacha20poly1305.Overhead
	if (payloadLen < 0) || (payloadLen > MaxPayloadLen) {
		return nil, uuid.UUID{}, errors.New("invalid packet length")
	}

	deviceId := uuid.UUID{}
	copy(deviceId[:], pkt[:16])
	nonce := pkt[16:16+12]
	ciphertext := pkt[HeaderLen:]
	payload, err := msg.Cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, uuid.UUID{}, err
	}

	return payload, uuid.UUID(deviceId), nil
}
