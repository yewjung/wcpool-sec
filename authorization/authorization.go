package authorization

import (
	context "context"
	"sec/controller"
	"sec/models"
)

type AuthorizationServerImpl struct {
	UnimplementedAuthorizationServer
	Storage models.Storage
}

func (auth AuthorizationServerImpl) VerifyPartyID(ctx context.Context, verification *Verification) (*VerificationResult, error) {
	token := verification.GetToken()
	authService := controller.AuthUserService{DB: auth.Storage.PostgresUserDB}
	ok, email := authService.IsTokenStillValid(token)
	if !ok {
		return &VerificationResult{
			Ok:    false,
			Email: email,
		}, nil
	}
	accountService := controller.AccountService{MongoDB: auth.Storage.MongoAccountDB, Cache: auth.Storage.RedisAccountCache}
	account := accountService.FindByEmail(email)

	verificationMethods := map[Option]models.VerificationMethod{
		Option_PARTY_ID: accountService.IsUserFromParty,
		Option_IS_ADMIN: accountService.IsUserAdminOfParty,
	}
	for _, option := range verification.GetOptions() {
		method, ok := verificationMethods[option]
		if !ok || !method(account, verification.Partyid) {
			return &VerificationResult{
				Ok:    false,
				Email: email,
			}, nil
		}
	}

	return &VerificationResult{
		Ok:    true,
		Email: email,
	}, nil

}
