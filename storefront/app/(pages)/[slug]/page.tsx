"use client";

import { useParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { publicApi } from "@/lib/api-client";
import { queryKeys } from "@/lib/query-keys";
import { t } from "@/lib/i18n";
import { Breadcrumb } from "@/components/layout/breadcrumb";

interface Page {
  id: string;
  title: string;
  slug: string;
  status: string;
}

interface Block {
  id: string;
  type: string;
  content: string;
  sort_order: number;
}

export default function CmsPage() {
  const params = useParams<{ slug: string }>();

  const { data, isLoading, error } = useQuery({
    queryKey: queryKeys.cms.page(params.slug),
    queryFn: async () => {
      const res = await publicApi.get<{ page: Page; blocks: Block[] }>(
        `/cms/pages/${params.slug}`
      );
      return res.data!;
    },
  });

  if (isLoading) {
    return <p className="text-muted-foreground">{t("cms.loading")}</p>;
  }

  if (error || !data) {
    return <p className="text-destructive">{t("cms.notFound")}</p>;
  }

  const { page, blocks } = data;
  const sorted = [...blocks].sort((a, b) => a.sort_order - b.sort_order);

  return (
    <article className="space-y-6">
      <Breadcrumb
        items={[
          { label: t("brand"), href: "/" },
          { label: page.title },
        ]}
      />
      <h1 className="text-3xl font-bold">{page.title}</h1>
      {sorted.map((block) => (
        <section key={block.id} className="cms-content leading-8">
          {block.type === "html" ? (
            <div dangerouslySetInnerHTML={{ __html: block.content }} />
          ) : (
            <p className="whitespace-pre-wrap">{block.content}</p>
          )}
        </section>
      ))}
      {sorted.length === 0 && (
        <p className="text-muted-foreground">{t("cms.empty")}</p>
      )}
    </article>
  );
}
