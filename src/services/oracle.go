package services

import (
	"fmt"
	"math/rand/v2"
)

func ConsultD6Oracle() string {
	dieThrow := rand.IntN(5) + 1
	switch dieThrow {
	case 1:
		return "No, and"
	case 2:
		return "No"
	case 3:
		return "No, but"
	case 4:
		return "Yes, but"
	case 5:
		return "Yes"
	case 6:
		return "Yes, and"
	default:
		return "No, and"
	}
}

func ConsultAdventureOracle() string {
	dieThrow := rand.IntN(100) + 1
	fmt.Println(dieThrow)
	if dieThrow <= 25 {
		return consultCaughtUpInEventsOracle()
	} else if dieThrow <= 50 {
		return consultEmploymentOracle()
	} else if dieThrow <= 75 {
		return consultExplorationOracle()
	} else {
		return consultCivicAffairOracle()
	}
}

func consultCaughtUpInEventsOracle() string {
	dieThrow := rand.IntN(100) + 1
	if dieThrow <= 16 {
		return fmt.Sprintf("Marked for death. I am wanted dead by: %s", consultWhoOracle())
	} else if dieThrow <= 33 {
		return fmt.Sprintf("Marked for death. I must help: %s. They are wanted dead by: %s", consultWhoOracle(), consultWhoOracle())
	} else if dieThrow <= 50 {
		return fmt.Sprintf("Blackmail. I am being blackmailed by: %s. I have to: %s", consultWhoOracle(), consultEmploymentOracle())
	} else if dieThrow <= 66 {
		return fmt.Sprintf("Escape. I've been captured and must escape or must help someone else escape capture. I may have to escape from imprisonment or I may need to help a kidnapping victim: %s. The opposition is: %s.", consultWhoOracle(), consultOppositionOracle())
	} else if dieThrow <= 83 {
		return fmt.Sprintf("A person in need. Someone runs up to me and begs for help. It could be a thief fleeing the wrath of their guild, or a victim fleeing a kidnapper: %s. The opposition is %s.", consultWhoOracle(), consultOppositionOracle())
	} else {
		return `Witness. You witness a deed that has serious consequences. This could be a kidnapping, a secret
				meeting, or a hunt for a fugitive. To find out what you see, roll on the Employment table. To discover
				your enemies, roll on the Opposition table`
	}
}

func consultEmploymentOracle() string {
	dieThrow := rand.IntN(100) + 1
	if dieThrow <= 9 {
		return fmt.Sprintf("Escort. I'm hired to escort a person or object to another location.")
	} else if dieThrow <= 18 {
		return fmt.Sprintf("Missing person. A person has gone missing and I've been asked to find them. The missing person is: %s. Behind the disappearance is: %s", consultWhoOracle(), consultOppositionOracle())
	} else if dieThrow <= 27 {
		return fmt.Sprintf("Harm (or Kill). I'm hired to harm/kill someone or something. My opposition is: %s", consultOppositionOracle())
	} else if dieThrow <= 36 {
		return fmt.Sprintf("Kidnap. I must kidnap: %s. My opposition is: %s", consultWhoOracle(), consultOppositionOracle())
	} else if dieThrow <= 45 {
		return fmt.Sprintf("Manhunt. A person is on the run and I'm hired to bring them back. The fugitive is: %s. They are being helped by: %s", consultWhoOracle(), consultOppositionOracle())
	} else if dieThrow <= 54 {
		return fmt.Sprintf("Guard duty. I'm hired by %s to protect against their enemies: %s", consultWhoOracle(), consultOppositionOracle())
	} else if dieThrow <= 63 {
		return fmt.Sprintf("Spying. I must gather information. My employer is: %s. I must spy on: %s", consultWhoOracle(), consultOppositionOracle())
	} else if dieThrow <= 72 {
		return fmt.Sprintf("Steal/Recover. I've been hired to steal from: %s", consultOppositionOracle())
	} else if dieThrow <= 81 {
		return fmt.Sprintf("Courier. I'm hired to deliver a package. My opposition is: %s", consultOppositionOracle())
	} else if dieThrow <= 90 {
		return fmt.Sprintf("Infiltrate. I must infiltrate an organization. My employer is: %s. My target is: %s. Target's influence level (1-6): %d", consultWhoOracle(), consultOppositionOracle(), rand.IntN(6)+1)
	} else {
		return fmt.Sprintf("Diplomacy. I'm hired to prevent escalation between factions. The aggrieved party is: %s", consultOppositionOracle())
	}
}

func consultExplorationOracle() string {
	dieThrow := rand.IntN(100) + 1
	if dieThrow <= 5 {
		return `Ruin`
	} else if dieThrow <= 10 {
		return `Mine`
	} else if dieThrow <= 15 {
		return `Cave`
	} else if dieThrow <= 20 {
		return `Subterranean vault`
	} else if dieThrow <= 25 {
		return `Forest`
	} else if dieThrow <= 30 {
		return `Desert`
	} else if dieThrow <= 35 {
		return `Stronghold`
	} else if dieThrow <= 40 {
		return `Fort`
	} else if dieThrow <= 45 {
		return `Small settlement`
	} else if dieThrow <= 50 {
		return `Tower`
	} else if dieThrow <= 55 {
		return `Swamp`
	} else if dieThrow <= 60 {
		return `Sewer`
	} else if dieThrow <= 65 {
		return `Graveyard`
	} else if dieThrow <= 70 {
		return `Temple/Church`
	} else if dieThrow <= 75 {
		return `Coastal area`
	} else if dieThrow <= 80 {
		return `Mansion`
	} else if dieThrow <= 85 {
		return `Catacomb`
	} else if dieThrow <= 90 {
		return `island`
	} else if dieThrow <= 95 {
		return `Laboratory`
	} else {
		return `Hideout`
	}
}

func consultCivicAffairOracle() string {
	dieThrow := rand.IntN(100) + 1
	if dieThrow <= 11 {
		return "Plague. I must find a cure for a disease that is spreading through an area."
	} else if dieThrow <= 22 {
		return "Natural disaster. I must deal with a disaster and aid those in need."
	} else if dieThrow <= 33 {
		return "Revolution. I'm involved in overthrowing a tyrannical leader or restoring the rightful ruler."
	} else if dieThrow <= 44 {
		return fmt.Sprintf("Criminal gangs. I need to drive out gangs terrorizing a community. Behind them is: %s", consultOppositionOracle())
	} else if dieThrow <= 55 {
		return "Mystical/technological threat. I must save a community from a technological or mystical threat."
	} else if dieThrow <= 66 {
		return "New religion/cult. I must investigate a new cult that has emerged in a town or city."
	} else if dieThrow <= 77 {
		return fmt.Sprintf("A valuable resource runs out. I must retrieve vital resources for a community. My opposition is: %s", consultOppositionOracle())
	} else if dieThrow <= 88 {
		return "Political upheaval. I've become involved in a major political upheaval."
	} else {
		return "Foreign threat. I must help a community threatened by hostile outsiders."
	}
}

func consultWhoOracle() string {
	dieThrow := rand.IntN(100) + 1
	if dieThrow <= 9 {
		return "A foreigner from a distant place"
	} else if dieThrow <= 18 {
		return "An ordinary person"
	} else if dieThrow <= 27 {
		return "A friend or lover"
	} else if dieThrow <= 36 {
		return "A representative of an organization"
	} else if dieThrow <= 45 {
		return "A family member"
	} else if dieThrow <= 54 {
		return "A mysterious figure"
	} else if dieThrow <= 63 {
		return "A religious figure"
	} else if dieThrow <= 72 {
		return fmt.Sprintf("A magic/tech user with influence level (1-6): %d", rand.IntN(6)+1)
	} else if dieThrow <= 81 {
		return fmt.Sprintf("An aristocrat with influence level (1-6): %d", rand.IntN(6)+1)
	} else if dieThrow <= 90 {
		return fmt.Sprintf("A government official with influence level (1-6): %d", rand.IntN(6)+1)
	} else {
		return fmt.Sprintf("A military person with rank level (1-6): %d", rand.IntN(6)+1)
	}
}

func consultOppositionOracle() string {
	dieThrow := rand.IntN(100) + 1
	if dieThrow <= 12 {
		return fmt.Sprintf("Cultists with strength level (1-6): %d", rand.IntN(6)+1)
	} else if dieThrow <= 25 {
		return fmt.Sprintf("Mercenaries with strength level (1-6): %d", rand.IntN(6)+1)
	} else if dieThrow <= 37 {
		return "Brigands/Gangs"
	} else if dieThrow <= 50 {
		return fmt.Sprintf("An aristocrat with influence level (1-6): %d", rand.IntN(6)+1)
	} else if dieThrow <= 62 {
		return fmt.Sprintf("An organization with influence level (1-6): %d", rand.IntN(6)+1)
	} else if dieThrow <= 75 {
		return "A magic/tech user"
	} else {
		return "An intelligent creature"
	}
}
