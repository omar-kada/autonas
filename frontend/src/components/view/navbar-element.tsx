import { cn } from '@/lib/utils';
import { Button } from '../ui/button';

export function NavbarElement(props: { label: string; navigate: () => void; className?: string }) {
  return (
    <Button
      variant="ghost"
      className={cn('text-sm font-medium', props.className)}
      onClick={() => props.navigate()}
    >
      {props.label}
    </Button>
  );
}
