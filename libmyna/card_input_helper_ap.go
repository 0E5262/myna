package libmyna

import (
	"errors"
	"fmt"
	"github.com/hamano/brokenasn1"
	"strconv"
)

type CardInputHelperAP struct {
	reader *Reader
}

type CardInputHelperAttrs struct {
	Header  []byte `asn1:"private,tag:33"`
	Name    string `asn1:"private,tag:34,utf8"`
	Address string `asn1:"private,tag:35,utf8"`
	Birth   string `asn1:"private,tag:36"`
	Sex     string `asn1:"private,tag:37"`
}

func (self *CardInputHelperAP) LookupPin() (int, error) {
	err := self.reader.SelectEF("00 11") // 券面事項入力補助用PIN
	if err != nil {
		return 0, err
	}
	count := self.reader.LookupPin()
	return count, nil
}

func (self *CardInputHelperAP) VerifyPin(pin string) error {
	err := self.reader.SelectEF("00 11")
	if err != nil {
		return err
	}
	err = self.reader.Verify(pin)
	return err
}

func (self *CardInputHelperAP) LookupPinA() (int, error) {
	err := self.reader.SelectEF("00 14")
	if err != nil {
		return 0, err
	}
	count := self.reader.LookupPin()
	return count, nil
}

func (self *CardInputHelperAP) VerifyPinA(pin string) error {
	err := self.reader.SelectEF("00 14")
	if err != nil {
		return err
	}
	err = self.reader.Verify(pin)
	return err
}

func (self *CardInputHelperAP) LookupPinB() (int, error) {
	err := self.reader.SelectEF("00 15")
	if err != nil {
		return 0, err
	}
	count := self.reader.LookupPin()
	return count, nil
}

func (self *CardInputHelperAP) VerifyPinB(pin string) error {
	err := self.reader.SelectEF("00 15")
	if err != nil {
		return err
	}
	err = self.reader.Verify(pin)
	return err
}

func (self *CardInputHelperAP) ReadMyNumber() (string, error) {
	err := self.reader.SelectEF("00 01")
	if err != nil {
		return "", err
	}
	data := self.reader.ReadBinary(16)
	var mynumber asn1.RawValue
	_, err = asn1.UnmarshalWithParams(data, &mynumber, "private,tag:16")
	if err != nil {
		return "", err
	}
	return string(mynumber.Bytes), nil
}

func (self *CardInputHelperAP) ReadAttributes() (*CardInputHelperAttrs, error) {
	err := self.reader.SelectEF("00 02")
	if err != nil {
		return nil, err
	}

	data := self.reader.ReadBinary(7)
	if len(data) != 7 {
		return nil, errors.New("Error at ReadBinary()")
	}

	parser := ASN1PartialParser{}
	err = parser.Parse(data)
	if err != nil {
		return nil, err
	}
	data = self.reader.ReadBinary(parser.GetSize())
	var attrs CardInputHelperAttrs
	_, err = asn1.UnmarshalWithParams(data, &attrs, "private,tag:32")
	if err != nil {
		return nil, err
	}
	return &attrs, nil
}

// ヘッダーをHEX文字列に変換
func (self *CardInputHelperAttrs) HeaderString() string {
	return fmt.Sprintf("% X", self.Header)
}

// ISO5218コードから日本語文字列に変換
func (self *CardInputHelperAttrs) SexString() string {
	n, err := strconv.Atoi(self.Sex)
	if err != nil {
		return "エラー"
	}
	switch n {
	case 1:
		return "男性"
	case 2:
		return "女性"
	case 9:
		return "適用不能"
	default:
		return "不明"
	}
}