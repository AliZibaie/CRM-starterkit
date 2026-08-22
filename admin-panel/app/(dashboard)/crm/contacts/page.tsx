"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { DataTable } from "@/components/table/data-table";

interface Contact {
  id: string;
  name: string;
  email: string;
  phone: string;
  company: string;
}

export default function ContactsPage() {
  const { data = [], isLoading } = useQuery({
    queryKey: queryKeys.contacts,
    queryFn: async () => {
      const res = await api.get<Contact[]>("/crm/contacts");
      return res.data ?? [];
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Contacts</h1>
        <p className="text-muted-foreground">CRM contacts directory</p>
      </div>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: "Name" },
          { key: "email", header: "Email" },
          { key: "phone", header: "Phone" },
          { key: "company", header: "Company" },
        ]}
      />
    </div>
  );
}
