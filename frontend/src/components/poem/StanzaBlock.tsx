import { useState } from 'react';
import { Link } from 'react-router-dom';
import type { Stanza } from '@/types/poem';
import { useToggleStanzaLike } from '@/hooks/useSocial';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';

interface StanzaBlockProps {
  stanza: Stanza;
  isLast?: boolean;
}

const deviceColors: Record<string, string> = {
  metaphor: 'bg-purple-100 text-purple-700',
  simile: 'bg-blue-100 text-blue-700',
  alliteration: 'bg-emerald-100 text-emerald-700',
  enjambment: 'bg-amber-100 text-amber-700',
  imagery: 'bg-rose-100 text-rose-700',
  personification: 'bg-cyan-100 text-cyan-700',
};

export function StanzaBlock({ stanza, isLast }: StanzaBlockProps) {
  const { isAuthenticated } = useAuthStore();
  const { openLogin } = useUIStore();
  const mutation = useToggleStanzaLike(stanza.poemId, stanza.id);
  const [liked, setLiked] = useState(stanza.likedByMe ?? false);
  const [count, setCount] = useState(stanza.likeCount);

  // Resync to server state after a refetch, adjusting state during render
  // rather than in an effect (see the poem like button).
  const [synced, setSynced] = useState({ liked: stanza.likedByMe ?? false, count: stanza.likeCount });
  if (!mutation.isPending && (synced.liked !== (stanza.likedByMe ?? false) || synced.count !== stanza.likeCount)) {
    setSynced({ liked: stanza.likedByMe ?? false, count: stanza.likeCount });
    setLiked(stanza.likedByMe ?? false);
    setCount(stanza.likeCount);
  }

  function toggleLike() {
    if (!isAuthenticated) {
      openLogin();
      return;
    }
    const next = !liked;
    setLiked(next);
    setCount((c) => c + (next ? 1 : -1));
    mutation.mutate(undefined, {
      onError: () => {
        setLiked(!next);
        setCount((c) => c + (next ? -1 : 1));
      },
    });
  }

  return (
    <div className="py-4">
      <div className="poem-text">{stanza.text}</div>

      <div className="mt-3 flex flex-wrap items-center gap-2">
        {stanza.author && (
          <Link
            to={`/profile/${stanza.authorId}`}
            className="font-sans text-sm text-feather no-underline hover:text-accent"
          >
            &mdash; @{stanza.author.displayName}
          </Link>
        )}

        {stanza.literaryDevice && (
          <span
            className={`inline-block rounded-full px-2 py-0.5 font-sans text-xs font-medium ${
              deviceColors[stanza.literaryDevice] ?? 'bg-gray-100 text-gray-700'
            }`}
          >
            {stanza.literaryDevice}
          </span>
        )}

        {stanza.status === 'pending' && (
          <span className="inline-block rounded-full bg-warning/15 px-2 py-0.5 font-sans text-xs font-medium text-warning">
            Awaiting approval
          </span>
        )}

        <button
          onClick={toggleLike}
          aria-label="Like stanza"
          className="ml-auto flex items-center gap-1 font-sans text-xs text-feather transition-colors hover:text-error"
        >
          <svg
            className={`h-4 w-4 transition-colors ${liked ? 'fill-error text-error' : 'fill-none'}`}
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path strokeLinecap="round" strokeLinejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
          {count > 0 && count}
        </button>
      </div>

      {!isLast && (
        <div className="mt-6 text-center text-feather/40 select-none">&#10022;</div>
      )}
    </div>
  );
}
