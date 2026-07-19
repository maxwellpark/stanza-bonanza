import { Link } from 'react-router-dom';
import { useDrafts } from '@/hooks/usePoems';
import { PoemCard } from '@/components/poem/PoemCard';

export function DraftsPage() {
  var { data, isLoading } = useDrafts();

  return (
    <div>
      <h1 className="mb-6 font-serif text-3xl font-bold text-ink">My Drafts</h1>

      {isLoading ? (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="card animate-pulse">
              <div className="mb-3 h-6 w-3/4 rounded bg-parchment-dark" />
              <div className="h-4 w-2/3 rounded bg-parchment-dark" />
            </div>
          ))}
        </div>
      ) : data && data.items.length > 0 ? (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {data.items.map((poem) => (
            <PoemCard key={poem.id} poem={poem} />
          ))}
        </div>
      ) : (
        <div className="card text-center text-feather">
          <p className="mb-4">No drafts yet.</p>
          <Link to="/poems/new" className="btn-primary">Start a poem</Link>
        </div>
      )}
    </div>
  );
}
