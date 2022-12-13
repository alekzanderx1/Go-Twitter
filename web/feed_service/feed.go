package feed

// Definition of Structs for Data storage

type Tweet struct {
	text             string
	createdBy        string
	createdTimestamp string
}

// In memory non-persistent storage

var posts = make(map[string]Tweet)

// getTweetsByUsers(userIds string)

// addNewTweet(Postdata)
