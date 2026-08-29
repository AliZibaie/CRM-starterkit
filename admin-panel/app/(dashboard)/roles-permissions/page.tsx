"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
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

  const { data: roles = [], isLoading: rolesLoading } = useQuery({
    queryKey: queryKeys.roles,
    queryFn: async () => {
      const res = await api.get<Role[]>("/rbac/roles");
      return res.data ?? [];
    },
  });

  const { data: permissions = [], isLoading: permsLoading } = useQuery({
    queryKey: queryKeys.permissions,
    queryFn: async () => {
      const res = await api.get<Permission[]>("/rbac/permissions");
      return res.data ?? [];
    },
  });

  const { data: roleDetail } = useQuery({
    queryKey: ["role", viewRoleId],
    queryFn: async () => {
      const res = await api.get<RoleDetail>(`/rbac/roles/${viewRoleId}`);
      return res.data!;
    },
    enabled: !!viewRoleId,
  });

  const createRole = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/roles", { name: roleName, description: roleDesc });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.roles });
      setRoleName("");
      setRoleDesc("");
    },
  });

  const createPerm = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/permissions", { key: permKey, name: permName });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.permissions });
      setPermKey("");
      setPermName("");
    },
  });

  const assignPerm = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/assign-permission", {
        role_id: selectedRole,
        permission_id: selectedPerm,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.roles });
      if (viewRoleId) {
        queryClient.invalidateQueries({ queryKey: ["role", viewRoleId] });
      }
    },
  });

  const assignRole = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/assign-role", {
        user_id: assignUserId,
        role_id: assignRoleId,
      });
    },
  });

  const syncPerms = useMutation({
    mutationFn: async () => {
      await api.post("/rbac/permissions/sync");
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: queryKeys.permissions }),
  });

  const rolePermissions = (roleDetail?.permission_ids ?? [])
    .map((id) => permissions.find((p) => p.id === id))
    .filter(Boolean) as Permission[];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("pages.roles.title")}</h1>
          <p className="text-muted-foreground">{t("pages.roles.subtitle")}</p>
        </div>
        <Button
          variant="outline"
          onClick={() => syncPerms.mutate()}
          disabled={syncPerms.isPending}
        >
          {t("pages.roles.syncPermissions")}
        </Button>
      </div>

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
            <Button onClick={() => createRole.mutate()} disabled={createRole.isPending}>
              {t("pages.roles.createRole")}
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
              <Input value={permKey} onChange={(e) => setPermKey(e.target.value)} />
            </div>
            <div className="space-y-2">
              <Label>{t("common.name")}</Label>
              <Input value={permName} onChange={(e) => setPermName(e.target.value)} />
            </div>
            <Button onClick={() => createPerm.mutate()} disabled={createPerm.isPending}>
              {t("pages.roles.createPermission")}
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
            <Button onClick={() => assignPerm.mutate()} disabled={assignPerm.isPending}>
              {t("common.assign")}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t("pages.roles.assignRoleToUser")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-2">
              <Label>{t("common.userId")}</Label>
              <Input
                value={assignUserId}
                onChange={(e) => setAssignUserId(e.target.value)}
              />
            </div>
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
            <Button onClick={() => assignRole.mutate()} disabled={assignRole.isPending}>
              {t("pages.roles.assignRole")}
            </Button>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("pages.roles.roleDetail")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <Select value={viewRoleId} onValueChange={setViewRoleId}>
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
            <div className="space-y-2">
              <p className="text-sm font-medium">{roleDetail.role.name}</p>
              <p className="text-sm text-muted-foreground">
                {roleDetail.role.description}
              </p>
              <div className="flex flex-wrap gap-2">
                {rolePermissions.map((p) => (
                  <span
                    key={p.id}
                    className="rounded-md border border-border px-2 py-1 text-xs"
                  >
                    {p.key}
                  </span>
                ))}
                {rolePermissions.length === 0 && (
                  <span className="text-sm text-muted-foreground">
                    {t("common.noPermissions")}
                  </span>
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
              { key: "key", header: t("common.key") },
              { key: "name", header: t("common.name") },
              { key: "module", header: t("common.module") },
            ]}
          />
        </div>
      </div>
    </div>
  );
}
