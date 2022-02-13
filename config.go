package main

import "os"

var TBToken = os.Getenv("TB")

// 删除骰子等 Dice Emoji
var DeleteDice = true

// 删除用户加入消息通知
var DeleteUserJoined = true

// 删除用户离开消息通知
var DeleteUserLeft = true

// 删除群组标题变更通知
var DeleteDNewGroupTitle = true

// 删除群组图片更换通知
var DeleteNewGroupPhoto = true

// 删除群组图片删除通知
var DeleteGroupPhotoDeleted = true

// 删除置顶通知
var DeleteOnPinned = true

// 删除频道消息
var DeleteChannel = true
