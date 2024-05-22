package cli

import (
	"fmt"
	"time"

	"github.com/nicolaa5/bankid"
)

func CollectRoutine(b bankid.BankID, orderRef string, output chan *bankid.CollectResponse, quit chan struct{}) {
	for {
		collectResponse, err := b.Collect(bankid.CollectRequest{
			OrderRef: orderRef,
		})
		if err != nil {
			fmt.Printf("Error collecting status: %v\n", err)
			close(output)
			quit <- struct{}{}
		}

		output <- collectResponse

		if collectResponse.Status == bankid.Pending {
			time.Sleep(1 * time.Second)
			continue
		}
		
		close(output)
		break
	}
}
