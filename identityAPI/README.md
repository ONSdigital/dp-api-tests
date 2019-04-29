
### IdentityAPI

I've archived the identityAPI tests for now as we're unlikely to be working on feature delivery any time soon.

To revert, decompress the included archive in the root of the dp-api-tests repo.

You'll also need to uncomment the following in testDataSetup/mongo/mongo.go
- line 4: `"github.com/ONSdigital/dp-api-tests/identityAPIModels"`
- lines 450-473 - the methods `GetIdentity` and `GetIdentities`

That will get you back to whatever test functionality he had pre-archiving.