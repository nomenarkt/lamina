import { z } from 'zod';

export const inviteUserSchema = z.object({
  email: z.string().email(),
  role: z.enum(['admin', 'viewer', 'external']),
  company: z.string().optional(),
  accessDuration: z
    .object({
      from: z.string().datetime(),
      to: z.string().datetime(),
    })
    .optional(),
});
