'use client';

import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { cn } from '@/lib';
import { ChevronDown } from 'lucide-react';
import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react';

interface AutoScrollAreaProps {
  children?: ReactNode;
  className?: string; // for the outer container
}

export function AutoScrollArea({ children, className = '' }: AutoScrollAreaProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  const [isBottomVisible, setIsBottomVisible] = useState(true);
  const [wasBottomVisible, setWasBottomVisible] = useState(true);
  const [userHasScrolled, setUserHasScrolled] = useState(false);

  // Helper to scroll smoothly to bottom
  const scrollToBottom = useCallback(() => {
    setUserHasScrolled(false);
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, []);

  // 1. Listen for explicit user scroll inputs (wheel, touch, key navigation)
  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const handleUserScroll = () => {
      setUserHasScrolled(true);
    };
    scrollToBottom();

    container.addEventListener('wheel', handleUserScroll, { passive: true });
    container.addEventListener('touchmove', handleUserScroll, { passive: true });
    container.addEventListener('keydown', handleUserScroll, { passive: true });

    return () => {
      container.removeEventListener('wheel', handleUserScroll);
      container.removeEventListener('touchmove', handleUserScroll);
      container.removeEventListener('keydown', handleUserScroll);
    };
  }, []);

  useEffect(() => {
    const newDataArrived = !isBottomVisible && wasBottomVisible && !userHasScrolled;
    if (newDataArrived) {
      scrollToBottom();
    } else if (!wasBottomVisible) {
      setWasBottomVisible(isBottomVisible);
    } else if (wasBottomVisible && isBottomVisible) {
      setUserHasScrolled(false);
    }
  }, [isBottomVisible, userHasScrolled, wasBottomVisible]);

  // Observe the bottom anchor with IntersectionObserver
  useEffect(() => {
    const bottomElement = bottomRef.current;
    if (!bottomElement) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        const isIntersecting = entry.isIntersecting;
        setIsBottomVisible(isIntersecting);
      },
      { root: null, threshold: 0.2 },
    );

    observer.observe(bottomElement);
    return () => observer.disconnect();
  }, []);

  return (
    <div ref={containerRef} className={cn('relative flex flex-col', className)}>
      {/* ScrollArea now uses h-full – it fills the parent */}
      <ScrollArea className="w-full h-full">
        <div className="p-4">
          {children}
          <div ref={bottomRef} className="h-1" />
        </div>
      </ScrollArea>

      {/* Floating button – fixed to viewport */}
      {userHasScrolled && (
        <Button
          onClick={scrollToBottom}
          className="fixed bottom-16 right-6 rounded-full p-3 shadow-lg z-50"
          size="icon"
        >
          <ChevronDown className="h-5 w-5" />
        </Button>
      )}
    </div>
  );
}
