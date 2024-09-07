package exports

// Routing Keys

/**
  - name: admin.registered
*/
var AdminRegisteredRoutingKeyName = "admin.registered"

/**
  - name: admin.send.welcome_email
*/
var AdminSendWelcomeEmailRoutingKeyName = "admin.send.welcome_email"

/**
  - name: admin_registered_by_admin.registered
*/
var AdminRegisteredByAdminRoutingKeyName = "admin_registered_by_admin.registered"

/**
  - name: admin_registered_by_admin.send.welcome_email
*/
var AdminRegisteredByAdminSendWelcomeEmailRoutingKeyName = "admin_registered_by_admin.send.welcome_email"

/**
  - name: judge.registered
*/
var JudgeRegisteredRoutingKeyName = "judge.registered"

/**
  - name: judge.send.welcome_email
*/
var JudgeSendWelcomeEmailRoutingKeyName = "judge.send.welcome_email"

/**
  - name: judge.registered
*/
var JudgeRegisteredByAdminRoutingKeyName = "judge.registered"

/**
  - name: judge_registered_by_admin.send.welcome_email
*/
var JudgeRegisteredByAdminSendWelcomeEmailRoutingKeyName = "judge_registered_by_admin.send.welcome_email"

/**
  - the exchange should be topic
  - name: participant.registered
*/
var ParticipantRegisteredRoutingKeyName = "participant.registered"

/**
  - the exchange should be topic
  - name: participant.send.welcome_email
*/
var ParticipantTeamLeadSendWelcomeEmailRoutingKeyName = "participant.team_lead.send.welcome_email"

// participant.send.welcome_email
var ParticipantTeamMemberWelcomeEmailQueueRoutingingKeyName = "participant.team_member.send.welcome_email"

/**
  - the exchange should be topic
  - name: participant.team_member.send.welcome_email
*/
var ParticipantTeamMemberSendWelcomeEmailRoutingKeyName = "participant.team_member.send.welcome_email"

/** ParticipantSendInvitationEmailRoutingKeyName
  - name: participant.invited
*/
var ParticipantInvitedRoutingKeyName = "participant.invited"

/**ParticipantSendInvitationEmailRoutingKeyName
  - name: participant.send.invitation_email
*/
var ParticipantSendInvitationEmailRoutingKeyName = "participant.send.invitation_email"

/*UploadJudgeProfilePicRoutingKeyName
  - judge.upload.profile_pic
*/
var UploadJudgeProfilePicRoutingKeyName = "judge.upload.profile_pic"
