import type { ContainerStatus } from '@/api/api';
import dockerLogo from '@/assets/docker.svg';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Item, ItemActions, ItemContent, ItemDescription, ItemMedia, ItemTitle } from '../ui/item';
import { HumanTime } from '../view/human-time';
import { ContainerStatusBadge } from './container-status-badge';

const logoPrefix = 'https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/png/';
export function ServiceStatus(props: {
  serviceName: string;
  serviceContainers: Array<ContainerStatus>;
}) {
  const time = props.serviceContainers[0]?.startedAt;
  return (
    <Item variant="outline">
      <ItemMedia>
        <Avatar className="size-10 rounded-none">
          <AvatarImage src={`${logoPrefix}${props.serviceName}.png`} />
          <AvatarFallback>
            <img src={dockerLogo} alt={`${props.serviceName} logo`} />
          </AvatarFallback>
        </Avatar>
      </ItemMedia>
      <ItemContent>
        <ItemTitle>{props.serviceName}</ItemTitle>
        <ItemDescription className="line-clamp-none">
          <HumanTime time={time} />
        </ItemDescription>
      </ItemContent>
      <ItemActions className="flex-wrap">
        {props.serviceContainers.map((item) => (
          <ContainerStatusBadge
            status={item.health}
            label={item.name}
            className="mx-1"
            key={`${props.serviceName}-${item.name}`}
          />
        ))}
      </ItemActions>
    </Item>
  );
}
