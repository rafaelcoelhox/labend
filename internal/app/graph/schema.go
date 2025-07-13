package graph

const Schema = `
	type User {
		id: ID!
		name: String!
		email: String!
		totalXP: Int!
		createdAt: String!
		updatedAt: String!
	}

	type UserXP {
		id: ID!
		userID: ID!
		sourceType: String!
		sourceID: String!
		amount: Int!
		createdAt: String!
	}

	type Challenge {
		id: ID!
		title: String!
		description: String!
		xpReward: Int!
		status: String!
		createdAt: String!
		updatedAt: String!
	}

	type ChallengeSubmission {
		id: ID!
		challengeID: ID!
		userID: ID!
		proofURL: String!
		status: String!
		createdAt: String!
	}

	type ChallengeVote {
		id: ID!
		submissionID: ID!
		userID: ID!
		approved: Boolean!
		timeCheck: Int!
		isValid: Boolean!
		createdAt: String!
	}

	input CreateUserInput {
		name: String!
		email: String!
	}

	input UpdateUserInput {
		name: String
		email: String
	}

	input CreateChallengeInput {
		title: String!
		description: String!
		xpReward: Int!
	}

	input SubmitChallengeInput {
		challengeID: ID!
		proofURL: String!
	}

	input VoteChallengeInput {
		submissionID: ID!
		approved: Boolean!
		timeCheck: Int!
	}

	type Query {
		user(id: ID!): User
		users(limit: Int = 10, offset: Int = 0): [User!]!
		userXPHistory(userID: ID!): [UserXP!]!
		
		challenge(id: ID!): Challenge
		challenges(limit: Int = 10, offset: Int = 0): [Challenge!]!
		challengeSubmissions(challengeID: ID!): [ChallengeSubmission!]!
		challengeVotes(submissionID: ID!): [ChallengeVote!]!
	}

	type Mutation {
		createUser(input: CreateUserInput!): User!
		updateUser(id: ID!, input: UpdateUserInput!): User!
		deleteUser(id: ID!): Boolean!
		
		createChallenge(input: CreateChallengeInput!): Challenge!
		submitChallenge(input: SubmitChallengeInput!): ChallengeSubmission!
		voteChallenge(input: VoteChallengeInput!): ChallengeVote!
	}
`
