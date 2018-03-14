package consensus

type Engine interface {
	VerifyHeader() error

	VerifyBlock() error

}
