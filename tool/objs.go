package tool
type Base struct {
    Id int	`col:"id" client"id"`	 // 显示顺序1
    BadgeId int	`col:"badgeId" client"badgeId"`	 // 徽章编号
    RuneType int	`col:"runeType" client"runeType"`	 // 可镶嵌符文类型类型
    RecycleReward IntSlice	`col:"recycleReward" client"recycleReward"`	 // 满星回收
    SkillId int	`col:"skillId" client"skillId"`	 // 普攻，对应skill
    RuneMax int	`col:"runeMax" client"runeMax"`	 // 符文等级上限
    LightMax int	`col:"lightMax" client"lightMax"`	 // 徽章升阶上限
    AddHp int	`col:"addHp" client"addHp"`	 // 自己给自己加血修正（百分比）
    BeAddHp int	`col:"beAddHp"`	 // 别人给自己加血修正（百分比）
}
