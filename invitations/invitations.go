package invitations

import (
	"net/mail"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Invitation struct {
	id               string
	to               string
	userID           string
	organizationName string
	acceptedAt       types.DateTime
}

func NewFromRecord(record *models.Record, dao *daos.Dao) (*Invitation, error) {
	orgID := record.GetStringDataValue("organization")
	collection, err := dao.FindCollectionByNameOrId("organizations")
	if err != nil {
		return nil, err
	}
	org, err := dao.FindRecordById(collection, orgID, nil)
	if err != nil {
		return nil, err
	}
	return &Invitation{
		id: record.Id,
		to: record.GetStringDataValue("to"),
		// acceptedAt:       record.GetDateTimeDataValue("accepted_at"),
		organizationName: org.GetStringDataValue("name"),
	}, nil
}

func SendInvitation(app core.App, rec *models.Record) error {
	client := app.NewMailClient()
	invitation, err := NewFromRecord(rec, app.Dao())
	if err != nil {
		return err
	}
	from, _ := mail.ParseAddress("vagmi@tarkalabs.com")
	to, _ := mail.ParseAddress(invitation.to)
	return client.Send(*from, *to, "Welcome on board", "welcome to tarkalabs. invite link", nil)
}
