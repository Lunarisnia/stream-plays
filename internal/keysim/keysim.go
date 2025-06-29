package keysim

import (
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
)

type KeySim interface {
	Press(key int)
}

type keySimImpl struct {
	keyBonding keybd_event.KeyBonding
}

func NewKeySim() (KeySim, error) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return nil, err
	}

	// For linux, it is very important to wait 2 seconds (WTF?)
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	return &keySimImpl{
		keyBonding: kb,
	}, nil
}

func (k *keySimImpl) Press(key int) {
	defer k.keyBonding.Clear()
	k.keyBonding.SetKeys(key)

	k.keyBonding.Press()
	time.Sleep(50 * time.Millisecond)
	k.keyBonding.Release()
}
