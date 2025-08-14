"use client";

import React, { Fragment, JSX } from "react";
import { usePathname } from "next/navigation";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbSeparator,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { AppRoute } from "@/lib/app-routes";

type TitleMap = {
  [K in AppRoute]?: string;
};

const TITLE_MAP: TitleMap = {
  [AppRoute.HOME]: "Home",
};

function humanize(segment: string): string {
  return segment.replace(/-/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

export const Breadcrumbs: React.FC = (): JSX.Element => {
  const pathname: string = usePathname();
  const segments: string[] = pathname.split("/").filter(Boolean);

  if (segments.length == 0) {
    segments.push(pathname);
  }

  return (
    <Breadcrumb>
      <BreadcrumbList>
        {segments.map((path, i) => {
          const isLast: boolean = i === segments.length - 1;
          const explicitTitle: string | undefined = TITLE_MAP[path as AppRoute];
          const dynamicSegment: string | undefined = path.split("/").pop()!;
          const title: string = explicitTitle ?? humanize(dynamicSegment);

          return (
            <Fragment key={path}>
              <BreadcrumbItem>
                {isLast ? (
                  <BreadcrumbPage>{title}</BreadcrumbPage>
                ) : (
                  <BreadcrumbLink href={path}>{title}</BreadcrumbLink>
                )}
              </BreadcrumbItem>
              {!isLast && <BreadcrumbSeparator />}
            </Fragment>
          );
        })}
      </BreadcrumbList>
    </Breadcrumb>
  );
};
