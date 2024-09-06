package exports

// Routing Keys

// the exchange should be topic
var JudgeRegisteredRoutingKeyName = "judge.registered"

// the exchange should be direct or topic
var JudgeSendWelcomeEmailRoutingKeyName = "judge.send.welcome_email"

/**
  - the exchange should be topic
  - name: participant.registered
*/
var ParticipantRegisteredRoutingKeyName = "participant.registered"

/**
  - the exchange should be topic
  - name: participant.send.welcome_email
*/
var ParticipantSendWelcomeEmailRoutingKeyName = "participant.send.welcome_email"

/** ParticipantSendInvitationEmailRoutingKeyName
  - name: participant.invited
*/
var ParticipantInvitedRoutingKeyName = "participant.invited"

/**ParticipantSendInvitationEmailRoutingKeyName
  - name: participant.send.invitation_email
*/
var ParticipantSendInvitationEmailRoutingKeyName = "participant.send.invitation_email"

/*UploadJudgeProfilePicBindingKeyName
  - judge.upload.profile_pic
*/
var UploadJudgeProfilePicRoutingKeyName = "judge.upload.profile_pic"
