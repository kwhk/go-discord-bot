package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func options(s *discordgo.Session, i *discordgo.InteractionCreate) {
	margs := []interface{}{
		i.ApplicationCommandData().Options[0].StringValue(),
		i.ApplicationCommandData().Options[1].IntValue(),
		i.ApplicationCommandData().Options[2].BoolValue(),
	}

	msgformat :=
		` Now you just learned how to use command options. Take a look to the value of which you've just entered:
		> string_option: %s
		> integer_option: %d
		> bool_option: %v `

	if len(i.ApplicationCommandData().Options) >= 4 {
		margs = append(margs, i.ApplicationCommandData().Options[3].ChannelValue(nil).ID)
		msgformat += "> channel-option: <#%s>\n"
	}
	if len(i.ApplicationCommandData().Options) >= 5 {
		margs = append(margs, i.ApplicationCommandData().Options[4].UserValue(nil).ID)
		msgformat += "> user-option: <@%s>\n"
	}
	if len(i.ApplicationCommandData().Options) >= 6 {
		margs = append(margs, i.ApplicationCommandData().Options[5].RoleValue(nil, "").ID)
		msgformat += "> role-option: <@&%s>\n"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, we'll discuss them in "responses" part
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}