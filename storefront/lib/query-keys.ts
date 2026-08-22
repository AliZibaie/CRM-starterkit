export const queryKeys = {
  products: {
    all: ["products"] as const,
    detail: (slug: string) => ["products", slug] as const,
  },
  categories: {
    all: ["categories"] as const,
    detail: (slug: string) => ["categories", slug] as const,
  },
  cart: ["cart"] as const,
  orders: {
    all: ["orders"] as const,
    detail: (id: string) => ["orders", id] as const,
  },
  tickets: ["tickets"] as const,
  cms: {
    page: (slug: string) => ["cms", "page", slug] as const,
    menus: ["cms", "menus"] as const,
  },
  payments: {
    gateways: ["payments", "gateways"] as const,
  },
};
