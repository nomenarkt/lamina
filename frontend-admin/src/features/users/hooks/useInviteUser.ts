import { useMutation } from '@tanstack/react-query';
import { inviteUser } from '../api';
import type { InviteUserPayload } from '../types';

export function useInviteUser() {
  return useMutation({
    mutationFn: (payload: InviteUserPayload) => inviteUser(payload),
  });
}
