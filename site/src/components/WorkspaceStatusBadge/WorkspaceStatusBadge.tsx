import Tooltip, {
  type TooltipProps,
  tooltipClasses,
} from "@mui/material/Tooltip";
import ErrorOutline from "@mui/icons-material/ErrorOutline";
import RecyclingIcon from "@mui/icons-material/Recycling";
import AutoDeleteIcon from "@mui/icons-material/AutoDelete";
import { type FC, type ReactNode } from "react";
import type { Workspace } from "api/typesGenerated";
import { Pill } from "components/Pill/Pill";
import { ChooseOne, Cond } from "components/Conditionals/ChooseOne";
import { DormantDeletionText } from "components/WorkspaceDeletion";
import { getDisplayWorkspaceStatus } from "utils/workspace";
import { useClassName } from "hooks/useClassName";
import { formatDistanceToNow } from "date-fns";

export type WorkspaceStatusBadgeProps = {
  workspace: Workspace;
  children?: ReactNode;
  className?: string;
};

export const WorkspaceStatusBadge: FC<WorkspaceStatusBadgeProps> = ({
  workspace,
  className,
}) => {
  const { text, icon, type } = getDisplayWorkspaceStatus(
    workspace.latest_build.status,
    workspace.latest_build.job,
  );

  return (
    <ChooseOne>
      <Cond condition={workspace.latest_build.status === "failed"}>
        <FailureTooltip
          title={
            <div css={{ display: "flex", alignItems: "center", gap: 10 }}>
              <ErrorOutline
                css={(theme) => ({
                  width: 14,
                  height: 14,
                  color: theme.palette.error.light,
                })}
              />
              <div>{workspace.latest_build.job.error}</div>
            </div>
          }
          placement="top"
        >
          <div>
            <Pill className={className} icon={icon} text={text} type={type} />
          </div>
        </FailureTooltip>
      </Cond>
      <Cond>
        <Pill className={className} icon={icon} text={text} type={type} />
      </Cond>
    </ChooseOne>
  );
};

export type DormantStatusBadgeProps = {
  workspace: Workspace;
  className?: string;
};

export const DormantStatusBadge: FC<DormantStatusBadgeProps> = ({
  workspace,
  className,
}) => {
  if (!workspace.dormant_at) {
    return <></>;
  }

  const formatDate = (dateStr: string): string => {
    const date = new Date(dateStr);
    return date.toLocaleDateString(undefined, {
      month: "long",
      day: "numeric",
      year: "numeric",
      hour: "numeric",
      minute: "numeric",
    });
  };

  return workspace.deleting_at ? (
    <Tooltip
      title={
        <>
          This workspace has not been used for{" "}
          {formatDistanceToNow(Date.parse(workspace.last_used_at))} and has been
          marked dormant. It is scheduled to be deleted on{" "}
          {formatDate(workspace.deleting_at)}.
        </>
      }
    >
      <Pill
        className={className}
        icon={<AutoDeleteIcon />}
        text="Deletion Pending"
        type="error"
      />
    </Tooltip>
  ) : (
    <Tooltip
      title={
        <>
          This workspace has not been used for{" "}
          {formatDistanceToNow(Date.parse(workspace.last_used_at))} and has been
          marked dormant. It is not scheduled for auto-deletion but will become
          a candidate if auto-deletion is enabled on this template.
        </>
      }
    >
      <Pill
        className={className}
        icon={<RecyclingIcon />}
        text="Dormant"
        type="warning"
      />
    </Tooltip>
  );
};

export const WorkspaceStatusText: FC<WorkspaceStatusBadgeProps> = ({
  workspace,
  className,
}) => {
  const { text, type } = getDisplayWorkspaceStatus(
    workspace.latest_build.status,
  );

  return (
    <ChooseOne>
      <Cond condition={Boolean(DormantDeletionText({ workspace }))}>
        <DormantDeletionText workspace={workspace} />
      </Cond>
      <Cond>
        <span
          role="status"
          data-testid="build-status"
          className={className}
          css={(theme) => ({
            fontWeight: 600,
            color: type
              ? theme.experimental.roles[type].fill
              : theme.experimental.l1.text,
          })}
        >
          {text}
        </span>
      </Cond>
    </ChooseOne>
  );
};

const FailureTooltip: FC<TooltipProps> = ({ children, ...tooltipProps }) => {
  const popper = useClassName(
    (css, theme) => css`
      & .${tooltipClasses.tooltip} {
        background-color: ${theme.palette.background.paper};
        border: 1px solid ${theme.palette.divider};
        font-size: 12px;
        padding: 8px 10px;
      }
    `,
    [],
  );

  return (
    <Tooltip {...tooltipProps} classes={{ popper }}>
      {children}
    </Tooltip>
  );
};
