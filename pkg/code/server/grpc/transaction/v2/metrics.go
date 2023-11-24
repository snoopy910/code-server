package transaction_v2

import (
	"context"
	"time"

	"github.com/code-payments/code-server/pkg/metrics"
	"github.com/code-payments/code-server/pkg/code/common"
	"github.com/code-payments/code-server/pkg/code/data/intent"
)

const (
	userIntentCreatedEventName            = "UserIntentCreated"
	privateUpgradeEventName               = "PrivateTransferUpgraded"
	submitIntentLatencyBreakdownEventName = "SubmitIntentLatencyBreakdown"
	airdropEventName                      = "Airdrop"
)

func recordUserIntentCreatedEvent(ctx context.Context, intentRecord *intent.Record) {
	metrics.RecordEvent(ctx, userIntentCreatedEventName, map[string]interface{}{
		"id":   intentRecord.IntentId,
		"type": intentRecord.IntentType.String(),
	})
}

func recordPrivacyUpgradedEvent(ctx context.Context, intentRecord *intent.Record, numUpgraded int) {
	upgradeTimeInMs := time.Since(intentRecord.CreatedAt) / time.Millisecond
	metrics.RecordEvent(ctx, privateUpgradeEventName, map[string]interface{}{
		"intent":             intentRecord.IntentId,
		"num_upgraded":       numUpgraded,
		"time_to_upgrade_ms": int(upgradeTimeInMs),
	})
}

func recordSubmitIntentLatencyBreakdownEvent(ctx context.Context, section string, latency time.Duration, actionCount int, intentType string) {
	latencyInMs := latency / time.Millisecond
	metrics.RecordEvent(ctx, submitIntentLatencyBreakdownEventName, map[string]interface{}{
		"section":      section,
		"latency_ms":   int(latencyInMs),
		"action_count": actionCount,
		"intent_type":  intentType,
	})
}

func recordAirdropEvent(ctx context.Context, owner *common.Account, airdropType AirdropType, usdValue float64) {
	metrics.RecordEvent(ctx, airdropEventName, map[string]interface{}{
		"owner":        owner.PublicKey().ToBase58(),
		"airdrop_type": airdropType.String(),
		"usd_value":    usdValue,
	})
}
