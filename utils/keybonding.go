package utils

import "github.com/micmonay/keybd_event"

func CopySelectionToClipboard(kb *keybd_event.KeyBonding) error {
	kb.HasCTRL(true)
	kb.HasSHIFT(false)
	kb.SetKeys(keybd_event.VK_C)
	err := kb.Launching()
	if err != nil {
		return err
	}
	return nil
}
