package mint

import (
	"fmt"

	sdk "github.com/ColorPlatform/color-sdk/types"
)

// Minter represents the minting state.
type Minter struct {
	Deflation        sdk.Dec `json:"deflation"`         // current annual inflation rate
	WeeklyProvisions sdk.Dec `json:"weekly_provisions"` // current weekly expected provisions
}

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(deflation, weeklyProvisions sdk.Dec) Minter {
	return Minter{
		Deflation:        deflation,
		WeeklyProvisions: weeklyProvisions,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(deflation sdk.Dec) Minter {
	return NewMinter(
		deflation,
		sdk.NewDec(300000000000),
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 13%.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(5, 2),
	)
}

func validateMinter(minter Minter) error {
	if minter.Deflation.LT(sdk.ZeroDec()) {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			minter.Deflation.String())
	}
	return nil
}

// NextWeeklySupply reduces the amount of weekly supply by 5%
func (m Minter) NextWeeklySupply() sdk.Dec {
	return m.WeeklyProvisions.Sub(m.Deflation.Mul(m.WeeklyProvisions))
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.WeeklyProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerWeek)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
