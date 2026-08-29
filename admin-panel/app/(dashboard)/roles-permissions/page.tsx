"use client";

import { useMemo, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-error-message";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Role {
  id: string;
  name: string;
  description: string;
}

interface Permission {
  id: string;
  key: string;
  name: string;
  module: string;
}

interface RoleDetail {
  role: Role;
  permission_ids: string[];
}

interface UserOption {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
}

export default function RolesPermissionsPage() {
  const queryClient = useQueryClient();
  const [roleName, setRoleName] = useState("");
  const [roleDesc, setRoleDesc] = useState("");
  const [permKey, setPermKey] = useState("");
  const [permName, setPermName] = useState("");
  const [selectedRole, setSelectedRole] = useState("");
  const [selectedPerm, setSelectedPerm] = useState("");
  const [assignUserId, setAssignUserId] = useState("");
  const [assignRoleId, setAssignRoleId] = useState("");
  const [viewRoleId, setViewRoleId] = useState("");
  const [togglingPermId, setTogglingPermId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ type: "error" | "success"; text: string } | null>(
    null
  );

  const { data: roles = [], isLoading: rolesLoading, isError: rolesError, error: rolesErr, refetch: refetchRoles } = useQuery({
    queryKey: queryKeys.roles,
    queryFn: async () => {
      const res = await api.get<Role[]>("/rbac/roles");
      return res.data ?? [];
    },
  });

  const { data: permissions = [], isLoading: permsLoading, isError: permsError, error: permsErr, refetch: refetchPerms } = useQuery({
    queryKey: queryKeys.permissions,
    queryFn: async () => {
      const res = await api.get<Permission[]>("/rbac/permissions");
      return res.data ?? [];
    },
  });

  const { data: users = [] } = useQuery({
    queryKey: ["users", "options"],
    queryFn: async () => {
      const res = await api.get<UserOption[]>("/users?page=1&per_page=100");
      return res.data ?? [];
    },
  });

  const { data: roleDetail, isLoading: roleDetailLoading } = useQuery({
    queryKey: ["role", viewRoleId],
    queryFn: async () => {
      const res = await api.get<RoleDetail>(`/rbac/roles/${viewRoleId}`);
      return res.data!;
    },
    enabled: !!viewRoleId,
  });

  const mutationOpts = {
    onError: (err: unknown) => setFeedback({ type: "error", text: getApiErrorMessage(err) }),
    onSuccess: () => setFeedback({ type: "success", text: t("common.operationSuccess") }),
  };

  const createRole = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/roles", { name: roleName, description: roleDesc });
    },
    onSuccess: () => {
      mutationOpts.onSuccess();
      queryClient.invalidateQueries({ queryKey: queryKeys.roles });
      setRoleName("");
      setRoleDesc("");
    },
    onError: mutationOpts.onError,
  });

  const createPerm = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/permissions", { key: permKey, name: permName });
    },
    onSuccess: () => {
      mutationOpts.onSuccess();
      queryClient.invalidateQueries({ queryKey: queryKeys.permissions });
      setPermKey("");
      setPermName("");
    },
    onError: mutationOpts.onError,
  });

  const assignPerm = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/assign-permission", {
        role_id: selectedRole,
        permission_id: selectedPerm,
      });
    },
    onSuccess: () => {
      mutationOpts.onSuccess();
      if (viewRoleId === selectedRole) {
        queryClient.invalidateQueries({ queryKey: ["role", viewRoleId] });
      }
    },
    onError: mutationOpts.onError,
  });

  const assignRole = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/assign-role", {
        user_id: assignUserId,
        role_id: assignRoleId,
      });
    },
    onSuccess: mutationOpts.onSuccess,
    onError: mutationOpts.onError,
  });

  const syncPerms = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/permissions/sync");
    },
    onSuccess: () => {
      mutationOpts.onSuccess();
      queryClient.invalidateQueries({ queryKey: queryKeys.permissions });
    },
    onError: mutationOpts.onError,
  });

  const togglePermission = async (permId: string, checked: boolean) => {
    if (!viewRoleId) return;
    setTogglingPermId(permId);
    setFeedback(null);
    try {
      if (checked) {
        await api.post("/rbac/assign-permission", {
          role_id: viewRoleId,
          permission_id: permId,
        });
      } else {
        await api.post("/rbac/revoke-permission", {
          role_id: viewRoleId,
          permission_id: permId,
        });
      }
      await queryClient.invalidateQueries({ queryKey: ["role", viewRoleId] });
      setFeedback({ type: "success", text: t("common.operationSuccess") });
    } catch (err) {
      setFeedback({ type: "error", text: getApiErrorMessage(err) });
    } finally {
      setTogglingPermId(null);
    }
  };

  const permissionsByModule = useMemo(() => {
    const groups: Record<string, Permission[]> = {};
    for (const p of permissions) {
      const mod = p.module || t("common.module");
      if (!groups[mod]) groups[mod] = [];
      groups[mod].push(p);
    }
    return groups;
  }, [permissions]);

  const assignedSet = new Set(roleDetail?.permission_ids ?? []);

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.roles.title")}</h1>
          <p className="text-muted-foreground">{t("pages.roles.subtitle")}</p>
        </div>
        <Button
          variant="outline"
          onClick={() => syncPerms.mutate()}
          disabled={syncPerms.isPending}
        >
          {syncPerms.isPending ? t("common.saving") : t("pages.roles.syncPermissions")}
        </Button>
      </div>

      {feedback && (
        <div
          className={`rounded-xl border px-4 py-3 text-sm ${
            feedback.type === "error"
              ? "border-destructive/40 bg-destructive/10 text-destructive"
              : "border-emerald-500/40 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
          }`}
        >
          {feedback.text}
        </div>
      )}

      {(rolesError || permsError) && (
        <div className="rounded-xl border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {getApiErrorMessage(rolesError ? rolesErr : permsErr)}
          <Button
            size="sm"
            variant="outline"
            className="ms-3"
            onClick={() => {
              if (rolesError) refetchRoles();
              if (permsError) refetchPerms();
            }}
          >
            {t("common.refresh")}
          </Button>
        </div>
      )}

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.roles.createRole")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.name")}</Label>
              <Input value={roleName} onChange={(e) => setRoleName(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>{t("common.description")}</Label>
              <Input value={roleDesc} onChange={(e) => setRoleDesc(e.target.value)} />
            </div>
            <Button
              onClick={() => createRole.mutate()}
              disabled={createRole.isPending || !roleName.trim()}
            >
              {createRole.isPending ? t("common.saving") : t("pages.roles.createRole")}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t("pages.roles.createPermission")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.key")}</Label>
              <Input value={permKey} onChange={(e) => setPermKey(e.target.value)} className="ltr-field" />
            </div>
            <div className="space-y-2">
              <Label>{t("common.name")}</Label>
              <Input value={permName} onChange={(e) => setPermName(e.target.value)} />
            </div>
            <Button
              onClick={() => createPerm.mutate()}
              disabled={createPerm.isPending || !permKey.trim() || !permName.trim()}
            >
              {createPerm.isPending ? t("common.saving") : t("pages.roles.createPermission")}
            </Button>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>{t("pages.roles.assignPermissions")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <Select value={selectedRole} onValueChange={setSelectedRole}>
              <SelectTrigger>
                <SelectValue placeholder={t("common.selectRole")} />
              </SelectTrigger>
              <SelectContent>
                {roles.map((r) => (
                  <SelectItem key={r.id} value={r.id}>
                    {r.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={selectedPerm} onValueChange={setSelectedPerm}>
              <SelectTrigger>
                <SelectValue placeholder={t("common.selectPermission")} />
              </SelectTrigger>
              <SelectContent>
                {permissions.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    {p.key}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              onClick={() => assignPerm.mutate()}
              disabled={assignPerm.isPending || !selectedRole || !selectedPerm}
            >
              {assignPerm.isPending ? t("common.saving") : t("common.assign")}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t("pages.roles.assignRoleToUser")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <Select value={assignUserId} onValueChange={setAssignUserId}>
              <SelectTrigger>
                <SelectValue placeholder={t("common.selectUser")} />
              </SelectTrigger>
              <SelectContent>
                {users.map((u) => (
                  <SelectItem key={u.id} value={u.id}>
                    {u.email}
                    {(u.first_name || u.last_name) &&
                      ` (${u.first_name} ${u.last_name})`.trim()}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={assignRoleId} onValueChange={setAssignRoleId}>
              <SelectTrigger>
                <SelectValue placeholder={t("common.selectRole")} />
              </SelectTrigger>
              <SelectContent>
                {roles.map((r) => (
                  <SelectItem key={r.id} value={r.id}>
                    {r.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              onClick={() => assignRole.mutate()}
              disabled={assignRole.isPending || !assignUserId || !assignRoleId}
            >
              {assignRole.isPending ? t("common.saving") : t("pages.roles.assignRole")}
            </Button>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.roles.roleDetail")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <Select
            value={viewRoleId}
            onValueChange={(id) => {
              setViewRoleId(id);
              setFeedback(null);
            }}
          >
            <SelectTrigger className="max-w-md">
              <SelectValue placeholder={t("common.selectRoleToView")} />
            </SelectTrigger>
            <SelectContent>
              {roles.map((r) => (
                <SelectItem key={r.id} value={r.id}>
                  {r.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {roleDetail && (
            <div className="space-y-4">
              <div>
                <p className="text-sm font-medium">{roleDetail.role.name}</p>
                <p className="text-sm text-muted-foreground">
                  {roleDetail.role.description || "—"}
                </p>
              </div>
              <div>
                <p className="mb-3 text-sm font-medium">{t("pages.roles.managePermissions")}</p>
                {roleDetailLoading ? (
                  <p className="text-sm text-muted-foreground">{t("common.loading")}</p>
                ) : (
                  <div className="max-h-80 space-y-4 overflow-y-auto rounded-xl border border-border p-4">
                    {Object.entries(permissionsByModule).map(([mod, perms]) => (
                      <div key={mod}>
                        <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                          {mod}
                        </p>
                        <div className="space-y-2">
                          {perms.map((p) => (
                            <label
                              key={p.id}
                              className="flex items-center gap-2 text-sm"
                            >
                              <input
                                type="checkbox"
                                checked={assignedSet.has(p.id)}
                                disabled={togglingPermId === p.id}
                                onChange={(e) => togglePermission(p.id, e.target.checked)}
                                className="h-4 w-4"
                              />
                              <span className="ltr-field font-mono text-xs">{p.key}</span>
                              <span className="text-muted-foreground">{p.name}</span>
                            </label>
                          ))}
                        </div>
                      </div>
                    ))}
                    {permissions.length === 0 && (
                      <p className="text-sm text-muted-foreground">{t("common.noPermissions")}</p>
                    )}
                  </div>
                )}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      <div className="grid gap-4 lg:grid-cols-2">
        <div>
          <h2 className="mb-3 text-lg font-semibold">{t("pages.roles.rolesList")}</h2>
          <DataTable
            loading={rolesLoading}
            data={roles as unknown as Record<string, unknown>[]}
            keyField="id"
            columns={[
              { key: "name", header: t("common.name") },
              { key: "description", header: t("common.description") },
              {
                key: "view",
                header: "",
                render: (row) => (
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={() => setViewRoleId(row.id as string)}
                  >
                    {t("common.view")}
                  </Button>
                ),
              },
            ]}
          />
        </div>
        <div>
          <h2 className="mb-3 text-lg font-semibold">{t("pages.roles.permissionsList")}</h2>
          <DataTable
            loading={permsLoading}
            data={permissions as unknown as Record<string, unknown>[]}
            keyField="id"
            columns={[
              { key: "key", header: t("common.key"), className: "ltr-field font-mono text-xs" },
              { key: "name", header: t("common.name") },
              { key: "module", header: t("common.module") },
            ]}
          />
        </div>
      </div>
    </div>
  );
}
