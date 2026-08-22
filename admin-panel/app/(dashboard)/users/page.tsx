"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DynamicForm, FormField } from "@/components/form-builder/dynamic-form";

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  is_active: boolean;
  created_at: string;
}

const createFields: FormField[] = [
  { name: "email", label: "Email", type: "email", required: true },
  { name: "password", label: "Password", type: "password", required: true },
  { name: "first_name", label: "First Name", type: "text" },
  { name: "last_name", label: "Last Name", type: "text" },
];

export default function UsersPage() {
  const queryClient = useQueryClient();
  const [page] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [formValues, setFormValues] = useState<Record<string, unknown>>({
    email: "",
    password: "",
    first_name: "",
    last_name: "",
  });
  const [editId, setEditId] = useState<string | null>(null);
  const [editValues, setEditValues] = useState<Record<string, unknown>>({});

  const { data, isLoading } = useQuery({
    queryKey: queryKeys.users(page),
    queryFn: async () => {
      const res = await api.get<User[]>(`/users?page=${page}&per_page=20`);
      return { users: res.data ?? [], meta: res.meta };
    },
  });

  const createMutation = useMutation({
    mutationFn: async (payload: Record<string, unknown>) => {
      await api.post("/users", payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setShowForm(false);
      setFormValues({ email: "", password: "", first_name: "", last_name: "" });
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({
      id,
      payload,
    }: {
      id: string;
      payload: Record<string, unknown>;
    }) => {
      await api.put(`/users/${id}`, payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      setEditId(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/users/${id}`);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["users"] }),
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Users</h1>
          <p className="text-muted-foreground">Manage system users</p>
        </div>
        <Button onClick={() => setShowForm(!showForm)}>
          {showForm ? "Cancel" : "Add User"}
        </Button>
      </div>

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>Create User</CardTitle>
          </CardHeader>
          <CardContent>
            <DynamicForm
              fields={createFields}
              values={formValues}
              onChange={(n, v) => setFormValues((p) => ({ ...p, [n]: v }))}
              onSubmit={() => createMutation.mutate(formValues)}
              submitLabel="Create"
              loading={createMutation.isPending}
            />
          </CardContent>
        </Card>
      )}

      <DataTable
        loading={isLoading}
        data={(data?.users ?? []) as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "email", header: "Email" },
          {
            key: "name",
            header: "Name",
            render: (row) =>
              `${row.first_name as string} ${row.last_name as string}`,
          },
          {
            key: "is_active",
            header: "Status",
            render: (row) => (
              <Badge variant={row.is_active ? "success" : "secondary"}>
                {row.is_active ? "Active" : "Inactive"}
              </Badge>
            ),
          },
          {
            key: "actions",
            header: "Actions",
            render: (row) => (
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => {
                    setEditId(row.id as string);
                    setEditValues({
                      email: row.email,
                      first_name: row.first_name,
                      last_name: row.last_name,
                      is_active: row.is_active,
                    });
                  }}
                >
                  Edit
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={() => deleteMutation.mutate(row.id as string)}
                >
                  Delete
                </Button>
              </div>
            ),
          },
        ]}
      />

      {editId && (
        <Card>
          <CardHeader>
            <CardTitle>Edit User</CardTitle>
          </CardHeader>
          <CardContent>
            <DynamicForm
              fields={[
                { name: "email", label: "Email", type: "email" },
                { name: "first_name", label: "First Name", type: "text" },
                { name: "last_name", label: "Last Name", type: "text" },
                { name: "is_active", label: "Active", type: "checkbox" },
              ]}
              values={editValues}
              onChange={(n, v) => setEditValues((p) => ({ ...p, [n]: v }))}
              onSubmit={() =>
                updateMutation.mutate({ id: editId, payload: editValues })
              }
              submitLabel="Save"
              loading={updateMutation.isPending}
            />
            <Button
              variant="ghost"
              className="mt-2"
              onClick={() => setEditId(null)}
            >
              Cancel
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
