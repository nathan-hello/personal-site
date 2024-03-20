import { db, Comment, NOW } from 'astro:db';

export default async function () {
    await db.insert(Comment).values([
        { parentId: 100018, body: 'Hope you like Astro DB!', created_at: NOW },
    ]);
}