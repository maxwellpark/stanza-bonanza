import { useState } from 'react';
import { Link } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';
import { NotificationsPanel } from './NotificationsPanel';
import { cn } from '@/lib/utils';

var publicLinks = [
  { to: '/explore', label: 'Explore' },
  { to: '/hall-of-fame', label: 'Hall of Fame' },
  { to: '/tutorials', label: 'Tutorials' },
];

export function Navbar() {
  var [mobileOpen, setMobileOpen] = useState(false);
  var { user, isAuthenticated } = useAuthStore();
  var { openLogin } = useUIStore();

  return (
    <nav className="sticky top-0 z-50 border-b border-parchment-dark bg-white/80 backdrop-blur">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <Link to="/" className="font-serif text-xl font-bold text-ink no-underline">
          Stanza Bonanza
        </Link>

        <div className="hidden items-center gap-6 md:flex">
          {publicLinks.map((link) => (
            <Link
              key={link.to}
              to={link.to}
              className="font-sans text-sm text-feather transition-colors hover:text-ink"
            >
              {link.label}
            </Link>
          ))}
          {isAuthenticated && (
            <Link to="/feed" className="font-sans text-sm text-feather transition-colors hover:text-ink">
              Feed
            </Link>
          )}
          {isAuthenticated && (
            <Link to="/drafts" className="font-sans text-sm text-feather transition-colors hover:text-ink">
              Drafts
            </Link>
          )}
        </div>

        <div className="hidden items-center gap-4 md:flex">
          {isAuthenticated && user ? (
            <div className="flex items-center gap-3">
              <NotificationsPanel />
              <Link to={`/profile/${user.id}`} className="flex items-center gap-2 no-underline">
                <img
                  src={user.avatarUrl || '/default-avatar.png'}
                  alt=""
                  className="h-8 w-8 rounded-full object-cover"
                />
                <span className="font-sans text-sm text-ink">{user.displayName}</span>
              </Link>
            </div>
          ) : (
            <button onClick={openLogin} className="btn-primary text-sm">
              Sign In
            </button>
          )}
        </div>

        <button
          className="flex min-h-[44px] min-w-[44px] flex-col items-center justify-center gap-1 md:hidden"
          onClick={() => setMobileOpen(!mobileOpen)}
          aria-label={mobileOpen ? 'Close menu' : 'Open menu'}
          aria-expanded={mobileOpen}
        >
          <span className={cn('block h-0.5 w-6 bg-ink transition-transform', mobileOpen && 'translate-y-1.5 rotate-45')} />
          <span className={cn('block h-0.5 w-6 bg-ink transition-opacity', mobileOpen && 'opacity-0')} />
          <span className={cn('block h-0.5 w-6 bg-ink transition-transform', mobileOpen && '-translate-y-1.5 -rotate-45')} />
        </button>
      </div>

      <AnimatePresence>
        {mobileOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden border-t border-parchment-dark md:hidden"
          >
            <div className="flex flex-col px-4 py-2">
              {publicLinks.map((link) => (
                <Link
                  key={link.to}
                  to={link.to}
                  onClick={() => setMobileOpen(false)}
                  className="flex min-h-[44px] items-center font-sans text-base text-feather transition-colors hover:text-ink"
                >
                  {link.label}
                </Link>
              ))}
              {isAuthenticated && (
                <Link
                  to="/feed"
                  onClick={() => setMobileOpen(false)}
                  className="flex min-h-[44px] items-center font-sans text-base text-feather transition-colors hover:text-ink"
                >
                  Feed
                </Link>
              )}
              {isAuthenticated && (
                <Link
                  to="/drafts"
                  onClick={() => setMobileOpen(false)}
                  className="flex min-h-[44px] items-center font-sans text-base text-feather transition-colors hover:text-ink"
                >
                  Drafts
                </Link>
              )}
              {isAuthenticated && user ? (
                <Link
                  to={`/profile/${user.id}`}
                  onClick={() => setMobileOpen(false)}
                  className="flex min-h-[44px] items-center gap-2 border-t border-parchment-dark pt-2 no-underline"
                >
                  <img
                    src={user.avatarUrl || '/default-avatar.png'}
                    alt=""
                    className="h-9 w-9 rounded-full object-cover"
                  />
                  <span className="font-sans text-base text-ink">{user.displayName}</span>
                </Link>
              ) : (
                <button
                  onClick={() => {
                    setMobileOpen(false);
                    openLogin();
                  }}
                  className="btn-primary mt-2"
                >
                  Sign In
                </button>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  );
}
