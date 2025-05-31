'use client';

import { useState } from 'react';
import { Tabs, Tab } from '@/components/ui/Tabs';
import UserRolesPanel from '@/features/access-control/UserRolesPanel';
import AccessPoliciesPanel from '@/features/access-control/AccessPoliciesPanel';

export default function AccessControlPage() {
  const [activeTab, setActiveTab] = useState<'roles' | 'policies'>('roles');

  return (
    <div className="p-6">
      <h1 className="text-xl font-semibold mb-4">Access Control</h1>
      <Tabs value={activeTab} onChange={setActiveTab}>
        <Tab label="User Roles" value="roles">
          <UserRolesPanel />
        </Tab>
        <Tab label="Access Policies" value="policies">
          <AccessPoliciesPanel />
        </Tab>
      </Tabs>
    </div>
  );
}
