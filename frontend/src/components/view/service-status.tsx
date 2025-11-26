import dockerLogo from '@/assets/docker.svg';
import type { ContainerStatus } from '@/models/stack-status';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Item, ItemContent, ItemDescription, ItemMedia, ItemTitle } from '../ui/item';
import { StatusBadge } from './status-badge';

const logoPrefix = 'https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/png/';
export function ServiceStatus(props: {
  serviceName: string;
  serviceContainers: Array<ContainerStatus>;
}) {
  return (
    <Item variant="outline">
      <ItemMedia>
        <Avatar className="size-10 rounded-none">
          <AvatarImage src={logoPrefix + props.serviceName + '.png'} />
          <AvatarFallback>
            <img src={dockerLogo} alt={`${props.serviceName} logo`} />
          </AvatarFallback>
        </Avatar>
      </ItemMedia>
      <ItemContent>
        <ItemTitle>{props.serviceName}</ItemTitle>
        <ItemDescription>
          {props.serviceContainers.map((item) => (
            <StatusBadge
              label={item.Name}
              status={item.Health}
              className="mx-1"
              key={`${props.serviceName}-${item.Name}`}
            />
          ))}
        </ItemDescription>
      </ItemContent>
    </Item>
  );
}
