export interface User {
  id: string;
  displayName: string;
  email?: string; // only present on /auth/me (the owner), not public profiles
  bio: string;
  avatarUrl: string;
  isVerified: boolean;
  createdAt: string;
  updatedAt: string;
}
