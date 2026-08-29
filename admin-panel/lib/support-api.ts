export interface TicketMessage {
  id: string;
  body: string;
  created_at: string;
  is_internal: boolean;
  sender_id?: string;
}

export interface TicketRecord {
  id: string;
  subject: string;
  status: string;
  priority: string;
  department_id?: string;
  assigned_to?: string;
  created_at: string;
  updated_at?: string;
}

export interface TicketDetailPayload {
  ticket: TicketRecord;
  messages: TicketMessage[];
}

export function parseTicketDetail(data: unknown): TicketDetailPayload {
  const raw = data as Record<string, unknown>;
  if (raw.ticket && raw.messages) {
    return {
      ticket: raw.ticket as TicketRecord,
      messages: raw.messages as TicketMessage[],
    };
  }
  const ticket = raw as unknown as TicketRecord;
  return {
    ticket,
    messages: (raw.messages as TicketMessage[] | undefined) ?? [],
  };
}
