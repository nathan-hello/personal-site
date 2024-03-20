import { defineDb, defineTable, column } from 'astro:db';

export const Comment = defineTable({
        columns: {
                id: column.number({ primaryKey: true }),
                parentId: column.number(),
                author: column.text({ default: "Anonymous" }),
                body: column.text(),
                created_at: column.date()
        }
});

export default defineDb({
        tables: { Comment },
});
