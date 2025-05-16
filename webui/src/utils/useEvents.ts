import { EventSource } from "eventsource";
import { useState, useEffect } from "react";


export type Event = {
  meta: {
    eventName: string;
    eventId: string;
    replyToId?: string;
    sourceModuleName?: string;
    time: string;
  };
} & Record<string, unknown>;

export type EventWrapper = {
  event: Event;
  replies: Event[];
};

export const useEvents = () => {
  const [events, setEvents] = useState<EventWrapper[]>([]);

  useEffect(() => {
    const es = new EventSource('http://localhost:20112/api/v0/events?target-id=*&target-name=*', {
      fetch: (input, init) => fetch(input, {
        ...init,
        headers: {
          ...init?.headers,
          'Last-Event-ID': init?.headers?.['Last-Event-ID'] ?? 'ALL_EVENTS',
          'Module-Name': 'debugger'
        }
      })
    });

    es.onerror = e => console.error(e);

    es.onmessage = e => {
      if (!e.data) return;

      try {
        const event = JSON.parse(e.data) as Event;

        if (event.meta.replyToId != undefined) {
          setEvents(events => {
            const result = [...events];
            const replies = result.find(r => r.event.meta.eventId === event.meta.replyToId)?.replies;
            if (replies && !replies.some(r => r.meta.eventId === event.meta.eventId)) replies?.push(event);
            return result;
          });
        } else {
          const internalEvent: EventWrapper = { event, replies: [] };

          setEvents(events => events.some(e => e.event.meta.eventId === event.meta.eventId)
            ? events
            : [internalEvent, ...events]);
        }
      } catch {
        console.error('failed to parse json:', e.data);
        return;
      }
    };

    return () => {
      es.close();
    };
  }, []);

  return events;
};