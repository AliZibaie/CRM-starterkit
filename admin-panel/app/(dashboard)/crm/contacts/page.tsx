"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
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
        <h1 className="text-2xl font-bold">{t("pages.contacts.title")}</h1>
        <p className="text-muted-foreground">{t("pages.contacts.subtitle")}</p>
      </div>

      <DataTable
        loading={isLoading}
        data={data as unknown as Record<string, unknown>[]}
        keyField="id"
        columns={[
          { key: "name", header: t("common.name") },
          { key: "email", header: t("common.email") },
          { key: "phone", header: t("common.phone") },
          { key: "company", header: t("common.company") },
        ]}
      />
    </div>
  );
}
