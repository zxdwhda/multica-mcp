# API coverage (0.4.0)

Source: https://github.com/multica-ai/multica/tree/7ebe0bf58d99238ccf830c7ed1aa658d4f5807d0

378 user-authenticated HTTP routes are registered as individual MCP tools. The original 16 convenience tools remain, plus one catalog query tool: 395 total.

This catalog records interface facts extracted from the official router and request declarations. It does not embed upstream server implementations. Live feature flags and Multica account permissions still apply. This is not a claim that all historical/future Multica versions support every route.

## Module coverage

| Module | Operations |
| --- | ---: |
| agent-activity-30d | 1 |
| agent-builder | 4 |
| agent-run-counts | 1 |
| agent-task-snapshot | 1 |
| agents | 25 |
| assignee-frequency | 1 |
| attachments | 4 |
| autopilots | 20 |
| chat | 25 |
| cli-token | 1 |
| client-usage | 1 |
| cloud-billing | 8 |
| cloud-runtime | 11 |
| cloud-subscriptions | 7 |
| comments | 9 |
| dashboard | 6 |
| dingtalk | 1 |
| feedback | 1 |
| inbox | 12 |
| integrations | 4 |
| invitations | 4 |
| issue-statuses | 5 |
| issue-view-preferences | 2 |
| issue-views | 5 |
| issues | 48 |
| labels | 5 |
| lark | 1 |
| me | 7 |
| notification-preferences | 3 |
| pins | 4 |
| plugin-bridge | 10 |
| projects | 10 |
| properties | 4 |
| quick-actions | 4 |
| runtimes | 17 |
| share-links | 1 |
| skills | 14 |
| slack | 1 |
| squads | 10 |
| tasks | 3 |
| telegram | 1 |
| tokens | 4 |
| upload-file | 1 |
| wecom | 1 |
| working-agents | 1 |
| workspaces | 69 |

## Operation coverage

| MCP tool | Method | API path |
| --- | --- | --- |
| `multica_api_get_workspace_agent_activity30d` | GET | `/api/agent-activity-30d` |
| `multica_api_save_agent_builder_draft` | PUT | `/api/agent-builder/sessions/{sessionId}/draft` |
| `multica_api_switch_agent_builder_runtime` | PATCH | `/api/agent-builder/sessions/{sessionId}/runtime` |
| `multica_api_list_agent_builder_sessions` | GET | `/api/agent-builder/sessions` |
| `multica_api_create_agent_builder_session` | POST | `/api/agent-builder/sessions` |
| `multica_api_get_workspace_agent_run_counts` | GET | `/api/agent-run-counts` |
| `multica_api_list_workspace_agent_task_snapshot` | GET | `/api/agent-task-snapshot` |
| `multica_api_create_mika_agent` | POST | `/api/agents/mika` |
| `multica_api_archive_agent` | POST | `/api/agents/{id}/archive` |
| `multica_api_cancel_agent_tasks` | POST | `/api/agents/{id}/cancel-tasks` |
| `multica_api_list_ding_talk_groups_for_agent` | GET | `/api/agents/{id}/dingtalk/groups` |
| `multica_api_get_agent_env` | GET | `/api/agents/{id}/env` |
| `multica_api_update_agent_env` | PUT | `/api/agents/{id}/env` |
| `multica_api_detach_label_from_agent` | DELETE | `/api/agents/{id}/labels/{labelId}` |
| `multica_api_list_labels_for_agent` | GET | `/api/agents/{id}/labels` |
| `multica_api_attach_label_to_agent` | POST | `/api/agents/{id}/labels` |
| `multica_api_set_agent_mcp_server_enabled` | PUT | `/api/agents/{id}/mcp-servers/{serverId}/enabled` |
| `multica_api_remove_agent_mcp_server` | DELETE | `/api/agents/{id}/mcp-servers/{serverId}` |
| `multica_api_list_agent_mcp_servers` | GET | `/api/agents/{id}/mcp-servers` |
| `multica_api_add_agent_mcp_server` | POST | `/api/agents/{id}/mcp-servers` |
| `multica_api_restore_agent` | POST | `/api/agents/{id}/restore` |
| `multica_api_set_agent_runtime_skill_enabled` | PUT | `/api/agents/{id}/runtime-skills/enabled` |
| `multica_api_add_agent_skills` | POST | `/api/agents/{id}/skills/add` |
| `multica_api_set_agent_skill_enabled` | PUT | `/api/agents/{id}/skills/{skillId}/enabled` |
| `multica_api_remove_agent_skill` | DELETE | `/api/agents/{id}/skills/{skillId}` |
| `multica_api_list_agent_skills` | GET | `/api/agents/{id}/skills` |
| `multica_api_set_agent_skills` | PUT | `/api/agents/{id}/skills` |
| `multica_api_list_agent_tasks` | GET | `/api/agents/{id}/tasks` |
| `multica_api_get_agent` | GET | `/api/agents/{id}` |
| `multica_api_update_agent` | PUT | `/api/agents/{id}` |
| `multica_api_list_agents` | GET | `/api/agents` |
| `multica_api_create_agent` | POST | `/api/agents` |
| `multica_api_get_assignee_frequency` | GET | `/api/assignee-frequency` |
| `multica_api_get_attachment_content` | GET | `/api/attachments/{id}/content` |
| `multica_api_download_attachment` | GET | `/api/attachments/{id}/download` |
| `multica_api_delete_attachment` | DELETE | `/api/attachments/{id}` |
| `multica_api_get_attachment_by_id` | GET | `/api/attachments/{id}` |
| `multica_api_cron_preview` | GET | `/api/autopilots/cron-preview` |
| `multica_api_get_autopilot_quota_usage` | GET | `/api/autopilots/usage` |
| `multica_api_remove_autopilot_collaborator` | DELETE | `/api/autopilots/{id}/collaborators/{userId}` |
| `multica_api_add_autopilot_collaborator` | POST | `/api/autopilots/{id}/collaborators` |
| `multica_api_replay_autopilot_delivery` | POST | `/api/autopilots/{id}/deliveries/{deliveryId}/replay` |
| `multica_api_get_autopilot_delivery` | GET | `/api/autopilots/{id}/deliveries/{deliveryId}` |
| `multica_api_list_autopilot_deliveries` | GET | `/api/autopilots/{id}/deliveries` |
| `multica_api_get_autopilot_run` | GET | `/api/autopilots/{id}/runs/{runId}` |
| `multica_api_list_autopilot_runs` | GET | `/api/autopilots/{id}/runs` |
| `multica_api_trigger_autopilot` | POST | `/api/autopilots/{id}/trigger` |
| `multica_api_rotate_autopilot_trigger_webhook_token` | POST | `/api/autopilots/{id}/triggers/{triggerId}/rotate-webhook-token` |
| `multica_api_set_autopilot_trigger_signing_secret` | PUT | `/api/autopilots/{id}/triggers/{triggerId}/signing-secret` |
| `multica_api_delete_autopilot_trigger` | DELETE | `/api/autopilots/{id}/triggers/{triggerId}` |
| `multica_api_update_autopilot_trigger` | PATCH | `/api/autopilots/{id}/triggers/{triggerId}` |
| `multica_api_create_autopilot_trigger` | POST | `/api/autopilots/{id}/triggers` |
| `multica_api_delete_autopilot` | DELETE | `/api/autopilots/{id}` |
| `multica_api_get_autopilot` | GET | `/api/autopilots/{id}` |
| `multica_api_update_autopilot` | PATCH | `/api/autopilots/{id}` |
| `multica_api_list_autopilots` | GET | `/api/autopilots` |
| `multica_api_create_autopilot` | POST | `/api/autopilots` |
| `multica_api_get_chat_channel_history` | GET | `/api/chat/history` |
| `multica_api_has_pending_chat_tasks` | GET | `/api/chat/pending-tasks/has-any` |
| `multica_api_list_pending_chat_tasks` | GET | `/api/chat/pending-tasks` |
| `multica_api_unpin_chat_agent` | DELETE | `/api/chat/pinned-agents/{agentId}` |
| `multica_api_list_chat_pinned_agents` | GET | `/api/chat/pinned-agents` |
| `multica_api_pin_chat_agent` | POST | `/api/chat/pinned-agents` |
| `multica_api_set_chat_session_archived` | PATCH | `/api/chat/sessions/{sessionId}/archive` |
| `multica_api_consume_chat_draft_restore` | DELETE | `/api/chat/sessions/{sessionId}/draft-restores/{restoreId}` |
| `multica_api_list_chat_draft_restores` | GET | `/api/chat/sessions/{sessionId}/draft-restores` |
| `multica_api_list_chat_messages_page` | GET | `/api/chat/sessions/{sessionId}/messages/page` |
| `multica_api_list_chat_messages` | GET | `/api/chat/sessions/{sessionId}/messages` |
| `multica_api_send_chat_message` | POST | `/api/chat/sessions/{sessionId}/messages` |
| `multica_api_start_mika_onboarding` | POST | `/api/chat/sessions/{sessionId}/onboarding` |
| `multica_api_get_pending_chat_task` | GET | `/api/chat/sessions/{sessionId}/pending-task` |
| `multica_api_set_chat_session_pinned` | PATCH | `/api/chat/sessions/{sessionId}/pin` |
| `multica_api_prioritize_queued_chat_task` | POST | `/api/chat/sessions/{sessionId}/queued-tasks/{taskId}/prioritize` |
| `multica_api_clear_queued_chat_tasks` | DELETE | `/api/chat/sessions/{sessionId}/queued-tasks` |
| `multica_api_regenerate_chat_quick_actions` | POST | `/api/chat/sessions/{sessionId}/quick-actions/regenerate` |
| `multica_api_mark_chat_session_read` | POST | `/api/chat/sessions/{sessionId}/read` |
| `multica_api_delete_chat_session` | DELETE | `/api/chat/sessions/{sessionId}` |
| `multica_api_get_chat_session` | GET | `/api/chat/sessions/{sessionId}` |
| `multica_api_update_chat_session` | PATCH | `/api/chat/sessions/{sessionId}` |
| `multica_api_list_chat_sessions` | GET | `/api/chat/sessions` |
| `multica_api_create_chat_session` | POST | `/api/chat/sessions` |
| `multica_api_get_chat_thread` | GET | `/api/chat/thread` |
| `multica_api_issue_cli_token` | POST | `/api/cli-token` |
| `multica_api_upsert_client_usage` | POST | `/api/client-usage` |
| `multica_api_get_cloud_billing_balance` | GET | `/api/cloud-billing/balance` |
| `multica_api_list_cloud_billing_batches` | GET | `/api/cloud-billing/batches` |
| `multica_api_get_cloud_billing_checkout_session` | GET | `/api/cloud-billing/checkout-sessions/{sessionId}` |
| `multica_api_create_cloud_billing_checkout_session` | POST | `/api/cloud-billing/checkout-sessions` |
| `multica_api_create_cloud_billing_portal_session` | POST | `/api/cloud-billing/portal-sessions` |
| `multica_api_list_cloud_billing_price_tiers` | GET | `/api/cloud-billing/price-tiers` |
| `multica_api_list_cloud_billing_topups` | GET | `/api/cloud-billing/topups` |
| `multica_api_list_cloud_billing_transactions` | GET | `/api/cloud-billing/transactions` |
| `multica_api_get_cloud_runtime_health` | GET | `/api/cloud-runtime/healthz` |
| `multica_api_exec_cloud_runtime_node` | POST | `/api/cloud-runtime/nodes/exec` |
| `multica_api_reboot_cloud_runtime_node` | POST | `/api/cloud-runtime/nodes/reboot` |
| `multica_api_start_cloud_runtime_node` | POST | `/api/cloud-runtime/nodes/start` |
| `multica_api_get_cloud_runtime_node_status` | POST | `/api/cloud-runtime/nodes/status` |
| `multica_api_stop_cloud_runtime_node` | POST | `/api/cloud-runtime/nodes/stop` |
| `multica_api_delete_cloud_runtime_node` | DELETE | `/api/cloud-runtime/nodes` |
| `multica_api_list_cloud_runtime_nodes` | GET | `/api/cloud-runtime/nodes` |
| `multica_api_create_cloud_runtime_node` | POST | `/api/cloud-runtime/nodes` |
| `multica_api_get_cloud_runtime_ready` | GET | `/api/cloud-runtime/readyz` |
| `multica_api_get_cloud_runtime_service` | GET | `/api/cloud-runtime` |
| `multica_api_create_cloud_workspace_subscription_checkout` | POST | `/api/cloud-subscriptions/checkout-sessions` |
| `multica_api_create_cloud_workspace_subscription_portal` | POST | `/api/cloud-subscriptions/portal-sessions` |
| `multica_api_get_cloud_workspace_subscription_prices` | GET | `/api/cloud-subscriptions/prices` |
| `multica_api_preview_cloud_workspace_subscription_seat_purchase` | POST | `/api/cloud-subscriptions/seats/purchase-preview` |
| `multica_api_purchase_cloud_workspace_subscription_seats` | POST | `/api/cloud-subscriptions/seats/purchases` |
| `multica_api_reconcile_cloud_workspace_subscription_seats` | POST | `/api/cloud-subscriptions/seats/reconcile` |
| `multica_api_get_cloud_workspace_subscription_summary` | GET | `/api/cloud-subscriptions/summary` |
| `multica_api_delete_comment_delete` | DELETE | `/api/comments/{commentId}/keep-replies` |
| `multica_api_remove_reaction` | DELETE | `/api/comments/{commentId}/reactions` |
| `multica_api_add_reaction` | POST | `/api/comments/{commentId}/reactions` |
| `multica_api_unresolve_comment` | DELETE | `/api/comments/{commentId}/resolve` |
| `multica_api_resolve_comment` | POST | `/api/comments/{commentId}/resolve` |
| `multica_api_preview_comment_sub_issue` | GET | `/api/comments/{commentId}/sub-issue-preview` |
| `multica_api_create_comment_sub_issue` | POST | `/api/comments/{commentId}/sub-issues` |
| `multica_api_delete_comment_delete_7cd56b06` | DELETE | `/api/comments/{commentId}` |
| `multica_api_update_comment` | PUT | `/api/comments/{commentId}` |
| `multica_api_get_dashboard_agent_run_time` | GET | `/api/dashboard/agent-runtime` |
| `multica_api_get_dashboard_failures_by_agent` | GET | `/api/dashboard/failures/by-agent` |
| `multica_api_get_dashboard_failures_daily` | GET | `/api/dashboard/failures/daily` |
| `multica_api_get_dashboard_run_time_daily` | GET | `/api/dashboard/runtime/daily` |
| `multica_api_get_dashboard_usage_by_agent` | GET | `/api/dashboard/usage/by-agent` |
| `multica_api_get_dashboard_usage_daily` | GET | `/api/dashboard/usage/daily` |
| `multica_api_redeem_ding_talk_binding_token` | POST | `/api/dingtalk/binding/redeem` |
| `multica_api_create_feedback` | POST | `/api/feedback` |
| `multica_api_archive_all_read_inbox` | POST | `/api/inbox/archive-all-read` |
| `multica_api_archive_all_inbox` | POST | `/api/inbox/archive-all` |
| `multica_api_archive_completed_inbox` | POST | `/api/inbox/archive-completed` |
| `multica_api_list_archived_inbox` | GET | `/api/inbox/archived` |
| `multica_api_mark_all_inbox_read` | POST | `/api/inbox/mark-all-read` |
| `multica_api_count_unread_inbox` | GET | `/api/inbox/unread-count` |
| `multica_api_unread_inbox_summary` | GET | `/api/inbox/unread-summary` |
| `multica_api_archive_inbox_item` | POST | `/api/inbox/{id}/archive` |
| `multica_api_mark_inbox_read` | POST | `/api/inbox/{id}/read` |
| `multica_api_unarchive_inbox_item` | POST | `/api/inbox/{id}/unarchive` |
| `multica_api_mark_inbox_unread` | POST | `/api/inbox/{id}/unread` |
| `multica_api_list_inbox` | GET | `/api/inbox` |
| `multica_api_composio_connect_init` | POST | `/api/integrations/composio/connect/init` |
| `multica_api_delete_composio_connection` | DELETE | `/api/integrations/composio/connections/{id}` |
| `multica_api_list_composio_connections` | GET | `/api/integrations/composio/connections` |
| `multica_api_list_composio_toolkits` | GET | `/api/integrations/composio/toolkits` |
| `multica_api_accept_invitation` | POST | `/api/invitations/{id}/accept` |
| `multica_api_decline_invitation` | POST | `/api/invitations/{id}/decline` |
| `multica_api_get_my_invitation` | GET | `/api/invitations/{id}` |
| `multica_api_list_my_invitations` | GET | `/api/invitations` |
| `multica_api_reorder_issue_statuses` | PATCH | `/api/issue-statuses/reorder` |
| `multica_api_archive_issue_status` | DELETE | `/api/issue-statuses/{id}` |
| `multica_api_update_issue_status` | PATCH | `/api/issue-statuses/{id}` |
| `multica_api_list_issue_statuses` | GET | `/api/issue-statuses` |
| `multica_api_create_issue_status` | POST | `/api/issue-statuses` |
| `multica_api_get_issue_view_preference` | GET | `/api/issue-view-preferences` |
| `multica_api_put_issue_view_preference` | PUT | `/api/issue-view-preferences` |
| `multica_api_delete_issue_view` | DELETE | `/api/issue-views/{id}` |
| `multica_api_get_issue_view_by_id` | GET | `/api/issue-views/{id}` |
| `multica_api_update_issue_view` | PATCH | `/api/issue-views/{id}` |
| `multica_api_list_issue_views` | GET | `/api/issue-views` |
| `multica_api_create_issue_view` | POST | `/api/issue-views` |
| `multica_api_batch_delete_issues` | POST | `/api/issues/batch-delete` |
| `multica_api_batch_update_issues` | POST | `/api/issues/batch-update` |
| `multica_api_child_issue_progress` | GET | `/api/issues/child-progress` |
| `multica_api_list_children_by_parents` | GET | `/api/issues/children` |
| `multica_api_list_grouped_issues` | GET | `/api/issues/grouped` |
| `multica_api_get_issue_limit_usage` | GET | `/api/issues/limit-usage` |
| `multica_api_preview_issue_trigger` | POST | `/api/issues/preview-trigger` |
| `multica_api_query_issues` | POST | `/api/issues/query` |
| `multica_api_quick_create_issue` | POST | `/api/issues/quick-create` |
| `multica_api_search_issues` | GET | `/api/issues/search` |
| `multica_api_list_issue_table_facets` | POST | `/api/issues/table/facets` |
| `multica_api_list_issue_table_groups` | POST | `/api/issues/table/groups` |
| `multica_api_list_issue_table_rows` | POST | `/api/issues/table/rows` |
| `multica_api_get_active_task_for_issue` | GET | `/api/issues/{id}/active-task` |
| `multica_api_list_attachments` | GET | `/api/issues/{id}/attachments` |
| `multica_api_list_child_issues` | GET | `/api/issues/{id}/children` |
| `multica_api_preview_comment_triggers` | POST | `/api/issues/{id}/comments/trigger-preview` |
| `multica_api_list_comments` | GET | `/api/issues/{id}/comments` |
| `multica_api_create_comment` | POST | `/api/issues/{id}/comments` |
| `multica_api_detach_label` | DELETE | `/api/issues/{id}/labels/{labelId}` |
| `multica_api_list_labels_for_issue` | GET | `/api/issues/{id}/labels` |
| `multica_api_attach_label` | POST | `/api/issues/{id}/labels` |
| `multica_api_delete_issue_metadata_key` | DELETE | `/api/issues/{id}/metadata/{key}` |
| `multica_api_set_issue_metadata_key` | PUT | `/api/issues/{id}/metadata/{key}` |
| `multica_api_list_issue_metadata` | GET | `/api/issues/{id}/metadata` |
| `multica_api_move_issue` | POST | `/api/issues/{id}/move` |
| `multica_api_delete_issue_property` | DELETE | `/api/issues/{id}/properties/{propertyId}` |
| `multica_api_set_issue_property` | PUT | `/api/issues/{id}/properties/{propertyId}` |
| `multica_api_list_pull_requests_for_issue` | GET | `/api/issues/{id}/pull-requests` |
| `multica_api_render_quick_action` | POST | `/api/issues/{id}/quick-actions/{quickActionId}/render` |
| `multica_api_run_quick_action` | POST | `/api/issues/{id}/quick-actions/{quickActionId}/run` |
| `multica_api_remove_issue_reaction` | DELETE | `/api/issues/{id}/reactions` |
| `multica_api_add_issue_reaction` | POST | `/api/issues/{id}/reactions` |
| `multica_api_rerun_issue` | POST | `/api/issues/{id}/rerun` |
| `multica_api_record_squad_leader_evaluation` | POST | `/api/issues/{id}/squad-evaluated` |
| `multica_api_subscribe_to_issue` | POST | `/api/issues/{id}/subscribe` |
| `multica_api_list_issue_subscribers` | GET | `/api/issues/{id}/subscribers` |
| `multica_api_list_tasks_by_issue` | GET | `/api/issues/{id}/task-runs` |
| `multica_api_cancel_task` | POST | `/api/issues/{id}/tasks/{taskId}/cancel` |
| `multica_api_list_timeline` | GET | `/api/issues/{id}/timeline` |
| `multica_api_unsubscribe_from_issue_subtree` | POST | `/api/issues/{id}/unsubscribe/subtree` |
| `multica_api_unsubscribe_from_issue` | POST | `/api/issues/{id}/unsubscribe` |
| `multica_api_get_issue_usage` | GET | `/api/issues/{id}/usage` |
| `multica_api_delete_issue` | DELETE | `/api/issues/{id}` |
| `multica_api_get_issue` | GET | `/api/issues/{id}` |
| `multica_api_update_issue` | PUT | `/api/issues/{id}` |
| `multica_api_list_issues` | GET | `/api/issues` |
| `multica_api_create_issue` | POST | `/api/issues` |
| `multica_api_delete_label` | DELETE | `/api/labels/{id}` |
| `multica_api_get_label` | GET | `/api/labels/{id}` |
| `multica_api_update_label` | PUT | `/api/labels/{id}` |
| `multica_api_list_labels` | GET | `/api/labels` |
| `multica_api_create_label` | POST | `/api/labels` |
| `multica_api_redeem_lark_binding_token` | POST | `/api/lark/binding/redeem` |
| `multica_api_join_cloud_waitlist` | POST | `/api/me/onboarding/cloud-waitlist` |
| `multica_api_complete_onboarding` | POST | `/api/me/onboarding/complete` |
| `multica_api_bootstrap_onboarding_no_runtime` | POST | `/api/me/onboarding/no-runtime-bootstrap` |
| `multica_api_bootstrap_onboarding_runtime` | POST | `/api/me/onboarding/runtime-bootstrap` |
| `multica_api_patch_onboarding` | PATCH | `/api/me/onboarding` |
| `multica_api_get_me` | GET | `/api/me` |
| `multica_api_update_me` | PATCH | `/api/me` |
| `multica_api_get_notification_preferences` | GET | `/api/notification-preferences` |
| `multica_api_patch_notification_preferences` | PATCH | `/api/notification-preferences` |
| `multica_api_update_notification_preferences` | PUT | `/api/notification-preferences` |
| `multica_api_reorder_pins` | PUT | `/api/pins/reorder` |
| `multica_api_delete_pin` | DELETE | `/api/pins/{itemType}/{itemId}` |
| `multica_api_list_pins` | GET | `/api/pins` |
| `multica_api_create_pin` | POST | `/api/pins` |
| `multica_api_get_plugin_context` | GET | `/api/plugin-bridge/v1/context` |
| `multica_api_invoke_plugin_hook` | POST | `/api/plugin-bridge/v1/hooks/{key}` |
| `multica_api_list_plugin_comments` | GET | `/api/plugin-bridge/v1/issues/{issue_ref}/comments` |
| `multica_api_create_plugin_comment` | POST | `/api/plugin-bridge/v1/issues/{issue_ref}/comments` |
| `multica_api_get_plugin_issue` | GET | `/api/plugin-bridge/v1/issues/{issue_ref}` |
| `multica_api_patch_plugin_issue` | PATCH | `/api/plugin-bridge/v1/issues/{issue_ref}` |
| `multica_api_delete_plugin_storage` | DELETE | `/api/plugin-bridge/v1/storage/{scope}/{key}` |
| `multica_api_get_plugin_storage` | GET | `/api/plugin-bridge/v1/storage/{scope}/{key}` |
| `multica_api_put_plugin_storage` | PUT | `/api/plugin-bridge/v1/storage/{scope}/{key}` |
| `multica_api_list_plugin_storage` | GET | `/api/plugin-bridge/v1/storage/{scope}` |
| `multica_api_search_projects` | GET | `/api/projects/search` |
| `multica_api_delete_project_resource` | DELETE | `/api/projects/{id}/resources/{resourceId}` |
| `multica_api_update_project_resource` | PUT | `/api/projects/{id}/resources/{resourceId}` |
| `multica_api_list_project_resources` | GET | `/api/projects/{id}/resources` |
| `multica_api_create_project_resource` | POST | `/api/projects/{id}/resources` |
| `multica_api_delete_project` | DELETE | `/api/projects/{id}` |
| `multica_api_get_project` | GET | `/api/projects/{id}` |
| `multica_api_update_project` | PUT | `/api/projects/{id}` |
| `multica_api_list_projects` | GET | `/api/projects` |
| `multica_api_create_project` | POST | `/api/projects` |
| `multica_api_get_property` | GET | `/api/properties/{id}` |
| `multica_api_update_property` | PATCH | `/api/properties/{id}` |
| `multica_api_list_properties` | GET | `/api/properties` |
| `multica_api_create_property` | POST | `/api/properties` |
| `multica_api_delete_quick_action` | DELETE | `/api/quick-actions/{id}` |
| `multica_api_update_quick_action` | PATCH | `/api/quick-actions/{id}` |
| `multica_api_list_quick_actions` | GET | `/api/quick-actions` |
| `multica_api_create_quick_action` | POST | `/api/quick-actions` |
| `multica_api_get_runtime_task_activity` | GET | `/api/runtimes/{runtimeId}/activity` |
| `multica_api_unbind_agents_and_delete_runtime_post` | POST | `/api/runtimes/{runtimeId}/archive-agents-and-delete` |
| `multica_api_get_local_skill_import_request` | GET | `/api/runtimes/{runtimeId}/local-skills/import/{requestId}` |
| `multica_api_initiate_import_local_skill` | POST | `/api/runtimes/{runtimeId}/local-skills/import` |
| `multica_api_get_local_skill_list_request` | GET | `/api/runtimes/{runtimeId}/local-skills/{requestId}` |
| `multica_api_initiate_list_local_skills` | POST | `/api/runtimes/{runtimeId}/local-skills` |
| `multica_api_get_model_list_request` | GET | `/api/runtimes/{runtimeId}/models/{requestId}` |
| `multica_api_initiate_list_models` | POST | `/api/runtimes/{runtimeId}/models` |
| `multica_api_unbind_agents_and_delete_runtime_post_4f3f3884` | POST | `/api/runtimes/{runtimeId}/unbind-agents-and-delete` |
| `multica_api_get_update` | GET | `/api/runtimes/{runtimeId}/update/{updateId}` |
| `multica_api_initiate_update` | POST | `/api/runtimes/{runtimeId}/update` |
| `multica_api_get_runtime_usage_by_agent` | GET | `/api/runtimes/{runtimeId}/usage/by-agent` |
| `multica_api_get_runtime_usage_by_hour` | GET | `/api/runtimes/{runtimeId}/usage/by-hour` |
| `multica_api_get_runtime_usage` | GET | `/api/runtimes/{runtimeId}/usage` |
| `multica_api_delete_agent_runtime` | DELETE | `/api/runtimes/{runtimeId}` |
| `multica_api_update_agent_runtime` | PATCH | `/api/runtimes/{runtimeId}` |
| `multica_api_list_agent_runtimes` | GET | `/api/runtimes` |
| `multica_api_join_by_share_link` | POST | `/api/share-links/join` |
| `multica_api_import_skill` | POST | `/api/skills/import` |
| `multica_api_search_skills` | GET | `/api/skills/search` |
| `multica_api_delete_skill_file` | DELETE | `/api/skills/{id}/files/{fileId}` |
| `multica_api_list_skill_files` | GET | `/api/skills/{id}/files` |
| `multica_api_upsert_skill_file` | PUT | `/api/skills/{id}/files` |
| `multica_api_detach_label_from_skill` | DELETE | `/api/skills/{id}/labels/{labelId}` |
| `multica_api_list_labels_for_skill` | GET | `/api/skills/{id}/labels` |
| `multica_api_attach_label_to_skill` | POST | `/api/skills/{id}/labels` |
| `multica_api_refresh_skill` | POST | `/api/skills/{id}/refresh` |
| `multica_api_delete_skill` | DELETE | `/api/skills/{id}` |
| `multica_api_get_skill` | GET | `/api/skills/{id}` |
| `multica_api_update_skill` | PUT | `/api/skills/{id}` |
| `multica_api_list_skills` | GET | `/api/skills` |
| `multica_api_create_skill` | POST | `/api/skills` |
| `multica_api_redeem_slack_binding_token` | POST | `/api/slack/binding/redeem` |
| `multica_api_update_squad_member_role` | PATCH | `/api/squads/{id}/members/role` |
| `multica_api_list_squad_member_status` | GET | `/api/squads/{id}/members/status` |
| `multica_api_remove_squad_member` | DELETE | `/api/squads/{id}/members` |
| `multica_api_list_squad_members` | GET | `/api/squads/{id}/members` |
| `multica_api_add_squad_member` | POST | `/api/squads/{id}/members` |
| `multica_api_delete_squad` | DELETE | `/api/squads/{id}` |
| `multica_api_get_squad` | GET | `/api/squads/{id}` |
| `multica_api_update_squad` | PUT | `/api/squads/{id}` |
| `multica_api_list_squads` | GET | `/api/squads` |
| `multica_api_create_squad` | POST | `/api/squads` |
| `multica_api_cancel_task_by_user` | POST | `/api/tasks/{taskId}/cancel` |
| `multica_api_list_task_messages_by_user` | GET | `/api/tasks/{taskId}/messages` |
| `multica_api_retry_source_context_quick_create` | POST | `/api/tasks/{taskId}/retry-source-context` |
| `multica_api_redeem_telegram_binding_token` | POST | `/api/telegram/binding/redeem` |
| `multica_api_renew_current_personal_access_token` | POST | `/api/tokens/current/renew` |
| `multica_api_revoke_personal_access_token` | DELETE | `/api/tokens/{id}` |
| `multica_api_list_personal_access_tokens` | GET | `/api/tokens` |
| `multica_api_create_personal_access_token` | POST | `/api/tokens` |
| `multica_api_upload_file` | POST | `/api/upload-file` |
| `multica_api_redeem_wecom_binding_token` | POST | `/api/wecom/binding/redeem` |
| `multica_api_list_workspace_working_agents` | GET | `/api/working-agents` |
| `multica_api_list_ding_talk_groups` | GET | `/api/workspaces/{id}/dingtalk/groups` |
| `multica_api_register_ding_talk_byo` | POST | `/api/workspaces/{id}/dingtalk/install/byo` |
| `multica_api_forget_ding_talk_group` | DELETE | `/api/workspaces/{id}/dingtalk/installations/{installationId}/groups/{conversationId}` |
| `multica_api_revoke_ding_talk_installation` | DELETE | `/api/workspaces/{id}/dingtalk/installations/{installationId}` |
| `multica_api_list_ding_talk_installations` | GET | `/api/workspaces/{id}/dingtalk/installations` |
| `multica_api_git_hub_connect` | GET | `/api/workspaces/{id}/github/connect` |
| `multica_api_list_git_hub_installation_repositories` | GET | `/api/workspaces/{id}/github/installations/{installationId}/repositories` |
| `multica_api_delete_git_hub_installation` | DELETE | `/api/workspaces/{id}/github/installations/{installationId}` |
| `multica_api_list_git_hub_installations` | GET | `/api/workspaces/{id}/github/installations` |
| `multica_api_revoke_invitation` | DELETE | `/api/workspaces/{id}/invitations/{invitationId}` |
| `multica_api_list_workspace_invitations` | GET | `/api/workspaces/{id}/invitations` |
| `multica_api_begin_lark_install` | POST | `/api/workspaces/{id}/lark/install/begin` |
| `multica_api_get_lark_install_status` | GET | `/api/workspaces/{id}/lark/install/{sessionId}/status` |
| `multica_api_revoke_lark_installation` | DELETE | `/api/workspaces/{id}/lark/installations/{installationId}` |
| `multica_api_list_lark_installations` | GET | `/api/workspaces/{id}/lark/installations` |
| `multica_api_leave_workspace` | POST | `/api/workspaces/{id}/leave` |
| `multica_api_delete_workspace_mcp_server` | DELETE | `/api/workspaces/{id}/mcp-servers/{serverId}` |
| `multica_api_update_workspace_mcp_server` | PUT | `/api/workspaces/{id}/mcp-servers/{serverId}` |
| `multica_api_list_workspace_mcp_servers` | GET | `/api/workspaces/{id}/mcp-servers` |
| `multica_api_create_workspace_mcp_server` | POST | `/api/workspaces/{id}/mcp-servers` |
| `multica_api_delete_member` | DELETE | `/api/workspaces/{id}/members/{memberId}` |
| `multica_api_update_member` | PATCH | `/api/workspaces/{id}/members/{memberId}` |
| `multica_api_list_members_with_user` | GET | `/api/workspaces/{id}/members` |
| `multica_api_create_invitation` | POST | `/api/workspaces/{id}/members` |
| `multica_api_publish_local_plugin_package` | POST | `/api/workspaces/{id}/plugins/packages/local` |
| `multica_api_delete_plugin_package` | DELETE | `/api/workspaces/{id}/plugins/packages/{packageId}` |
| `multica_api_list_plugin_packages` | GET | `/api/workspaces/{id}/plugins/packages` |
| `multica_api_publish_plugin_package` | POST | `/api/workspaces/{id}/plugins/packages` |
| `multica_api_preview_plugin` | POST | `/api/workspaces/{id}/plugins/preview` |
| `multica_api_configure_plugin` | PUT | `/api/workspaces/{id}/plugins/{installationId}/config` |
| `multica_api_disable_plugin` | POST | `/api/workspaces/{id}/plugins/{installationId}/disable` |
| `multica_api_enable_plugin` | POST | `/api/workspaces/{id}/plugins/{installationId}/enable` |
| `multica_api_list_plugin_invocations` | GET | `/api/workspaces/{id}/plugins/{installationId}/invocations` |
| `multica_api_list_plugin_mcp_tools` | GET | `/api/workspaces/{id}/plugins/{installationId}/mcp/{hookKey}/tools` |
| `multica_api_approve_plugin_mcp_tools` | PUT | `/api/workspaces/{id}/plugins/{installationId}/mcp/{hookKey}/tools` |
| `multica_api_get_plugin_surface_launch` | GET | `/api/workspaces/{id}/plugins/{installationId}/surfaces/{surfaceKey}/launch` |
| `multica_api_revoke_plugin_token` | DELETE | `/api/workspaces/{id}/plugins/{installationId}/token` |
| `multica_api_rotate_plugin_token` | POST | `/api/workspaces/{id}/plugins/{installationId}/token` |
| `multica_api_uninstall_plugin` | DELETE | `/api/workspaces/{id}/plugins/{installationId}` |
| `multica_api_list_plugins` | GET | `/api/workspaces/{id}/plugins` |
| `multica_api_install_plugin` | POST | `/api/workspaces/{id}/plugins` |
| `multica_api_delete_runtime_profile` | DELETE | `/api/workspaces/{id}/runtime-profiles/{profileId}` |
| `multica_api_get_runtime_profile` | GET | `/api/workspaces/{id}/runtime-profiles/{profileId}` |
| `multica_api_update_runtime_profile_patch` | PATCH | `/api/workspaces/{id}/runtime-profiles/{profileId}` |
| `multica_api_update_runtime_profile_put` | PUT | `/api/workspaces/{id}/runtime-profiles/{profileId}` |
| `multica_api_list_runtime_profiles` | GET | `/api/workspaces/{id}/runtime-profiles` |
| `multica_api_create_runtime_profile` | POST | `/api/workspaces/{id}/runtime-profiles` |
| `multica_api_revoke_share_link` | DELETE | `/api/workspaces/{id}/share-links/{linkId}` |
| `multica_api_list_share_links` | GET | `/api/workspaces/{id}/share-links` |
| `multica_api_create_share_link` | POST | `/api/workspaces/{id}/share-links` |
| `multica_api_register_slack_byo` | POST | `/api/workspaces/{id}/slack/install/byo` |
| `multica_api_revoke_slack_installation` | DELETE | `/api/workspaces/{id}/slack/installations/{installationId}` |
| `multica_api_list_slack_installations` | GET | `/api/workspaces/{id}/slack/installations` |
| `multica_api_register_telegram_bot` | POST | `/api/workspaces/{id}/telegram/install` |
| `multica_api_revoke_telegram_installation` | DELETE | `/api/workspaces/{id}/telegram/installations/{installationId}` |
| `multica_api_list_telegram_installations` | GET | `/api/workspaces/{id}/telegram/installations` |
| `multica_api_rotate_vcs_connection_webhook` | POST | `/api/workspaces/{id}/vcs/connections/{connectionId}/rotate-webhook` |
| `multica_api_delete_vcs_connection` | DELETE | `/api/workspaces/{id}/vcs/connections/{connectionId}` |
| `multica_api_list_vcs_connections` | GET | `/api/workspaces/{id}/vcs/connections` |
| `multica_api_connect_vcs` | POST | `/api/workspaces/{id}/vcs/connections` |
| `multica_api_register_wecom_byo` | POST | `/api/workspaces/{id}/wecom/install/byo` |
| `multica_api_revoke_wecom_installation` | DELETE | `/api/workspaces/{id}/wecom/installations/{installationId}` |
| `multica_api_list_wecom_installations` | GET | `/api/workspaces/{id}/wecom/installations` |
| `multica_api_delete_workspace` | DELETE | `/api/workspaces/{id}` |
| `multica_api_get_workspace` | GET | `/api/workspaces/{id}` |
| `multica_api_update_workspace_patch` | PATCH | `/api/workspaces/{id}` |
| `multica_api_update_workspace_put` | PUT | `/api/workspaces/{id}` |
| `multica_api_list_workspaces` | GET | `/api/workspaces` |
| `multica_api_create_workspace` | POST | `/api/workspaces` |

## Other authentication surfaces

The following routes use daemon, public callback/webhook, or Plugin installation-token authentication instead of the configured user PAT. They are listed for scope clarity, not silently represented as callable user tools. The user-authenticated Plugin bridge is included above. WebSocket channels and health endpoints are infrastructure, not request/response business tools.

| Method | Path | Surface |
| --- | --- | --- |
| GET | `/health` | Non-user API surface: public |
| GET | `/readyz` | Non-user API surface: public |
| GET | `/healthz` | Non-user API surface: public |
| GET | `/uploads/*` | Non-user API surface: public |
| GET | `/api/attachments/{id}/signed-download` | Non-user API surface: public |
| GET | `/api/avatars/{sig}/*` | Non-user API surface: public |
| GET | `/plugin-surfaces/{token}` | Non-user API surface: public |
| POST | `/auth/send-code` | Non-user API surface: public |
| POST | `/auth/verify-code` | Non-user API surface: public |
| POST | `/auth/google` | Non-user API surface: public |
| POST | `/auth/logout` | Non-user API surface: public |
| GET | `/api/config` | Non-user API surface: public |
| POST | `/api/contact-sales` | Non-user API surface: public |
| GET | `/api/share-links/{code}` | Non-user API surface: public |
| POST | `/api/webhooks/autopilots/{token}` | Non-user API surface: public |
| POST | `/api/webhooks/github` | Non-user API surface: public |
| GET | `/api/github/setup` | Non-user API surface: public |
| POST | `/api/webhooks/vcs/{connectionId}` | Non-user API surface: public |
| POST | `/api/webhooks/stripe` | Non-user API surface: public |
| GET | `/api/integrations/composio/callback` | Non-user API surface: public |
| POST | `/api/daemon/register` | Non-user API surface: daemon |
| POST | `/api/daemon/deregister` | Non-user API surface: daemon |
| POST | `/api/daemon/heartbeat` | Non-user API surface: daemon |
| GET | `/api/daemon/ws` | Non-user API surface: daemon |
| GET | `/api/daemon/workspaces` | Non-user API surface: daemon |
| GET | `/api/daemon/workspaces/{workspaceId}/repos` | Non-user API surface: daemon |
| GET | `/api/daemon/workspaces/{workspaceId}/runtime-profiles` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{id}/plugin-hooks` | Non-user API surface: daemon |
| GET | `/api/daemon/tasks/{id}/plugin-mcp/{contributionId}/credential` | Non-user API surface: daemon |
| POST | `/api/daemon/runtimes/{runtimeId}/tasks/claim` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/claim` | Non-user API surface: daemon |
| POST | `/api/daemon/claim` | Non-user API surface: daemon |
| POST | `/api/daemon/runtimes/{runtimeId}/tasks/{taskId}/prepare-lease` | Non-user API surface: daemon |
| POST | `/api/daemon/runtimes/{runtimeId}/tasks/{taskId}/skill-bundles/resolve` | Non-user API surface: daemon |
| GET | `/api/daemon/runtimes/{runtimeId}/tasks/pending` | Non-user API surface: daemon |
| POST | `/api/daemon/runtimes/{runtimeId}/update/{updateId}/result` | Non-user API surface: daemon |
| POST | `/api/daemon/runtimes/{runtimeId}/models/{requestId}/result` | Non-user API surface: daemon |
| POST | `/api/daemon/runtimes/{runtimeId}/local-skills/{requestId}/result` | Non-user API surface: daemon |
| POST | `/api/daemon/runtimes/{runtimeId}/local-skills/import/{requestId}/result` | Non-user API surface: daemon |
| GET | `/api/daemon/tasks/{taskId}/status` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/start` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/wait-local-directory` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/progress` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/complete` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/fail` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/usage` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/messages` | Non-user API surface: daemon |
| GET | `/api/daemon/tasks/{taskId}/messages` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/cancel-ack` | Non-user API surface: daemon |
| POST | `/api/daemon/workspaces/{workspaceId}/issues/gc-check` | Non-user API surface: daemon |
| GET | `/api/daemon/issues/{issueId}/gc-check` | Non-user API surface: daemon |
| GET | `/api/daemon/chat-sessions/{sessionId}/gc-check` | Non-user API surface: daemon |
| GET | `/api/daemon/autopilot-runs/{runId}/gc-check` | Non-user API surface: daemon |
| GET | `/api/daemon/tasks/{taskId}/gc-check` | Non-user API surface: daemon |
| POST | `/api/daemon/runtimes/{runtimeId}/recover-orphans` | Non-user API surface: daemon |
| POST | `/api/daemon/tasks/{taskId}/session` | Non-user API surface: daemon |
| GET | `/v1/context` | Non-user API surface: plugin |
| GET | `/v1/issues/{issue_ref}` | Non-user API surface: plugin |
| PATCH | `/v1/issues/{issue_ref}` | Non-user API surface: plugin |
| GET | `/v1/issues/{issue_ref}/comments` | Non-user API surface: plugin |
| POST | `/v1/issues/{issue_ref}/comments` | Non-user API surface: plugin |
| GET | `/v1/storage/{scope}` | Non-user API surface: plugin |
| GET | `/v1/storage/{scope}/{key}` | Non-user API surface: plugin |
| PUT | `/v1/storage/{scope}/{key}` | Non-user API surface: plugin |
| DELETE | `/v1/storage/{scope}/{key}` | Non-user API surface: plugin |

## Verification layers

- Every generated tool is registered and dispatched through the MCP SDK in automated tests.
- HTTP transport tests cover request paths, queries, empty and null fields, credentials, redirects, multipart upload, 204 responses, and complete upstream error envelopes.
- Real cloud API and ChatGPT acceptance results are recorded separately in verification-full-api.json after live tests.
- Membership, token revocation, billing, external integration installation, and execution-start operations are exposed but are not bulk-executed as a test.

## Regeneration

```sh
go run ./scripts/generate_api /path/to/multica-source SOURCE_COMMIT internal/apicatalog/catalog.json
go test -race ./...
```

Inputs use path for route identifiers, query for URL parameters, body for JSON, and files for base64 multipart data. Additional JSON body fields are forwarded for compatibility. Uploads are bounded to 4 MiB per file, 8 MiB per upstream request; responses are bounded to 8 MiB. Pagination remains the upstream contract; use its cursors/offsets and inspect the returned body.
