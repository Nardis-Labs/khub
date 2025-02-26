import React from "react";
import { ResourceDataTable } from "../../components/ResourceDataTable/ResourceDataTable";
import { Identification, Events, GroupSecurity, Misuse, CheckmarkFilled, Edit, Close } from "@carbon/icons-react";
import { useGetGroupsQuery, useGetPermissionsQuery, useGetUsersQuery, useUpsertGroupMutation, useUpsertPermissionMutation } from "../../../service/khub";
import { AdminDataTable } from "../../components/AdminDataTable/AdminDataTable";
import { useSelector } from "react-redux";
import { RootState, useAppDispatch } from "../../store";
import { updateAdminModalState } from "../../../service/adminmodals";


import { Tag, Button, ButtonSet, ComposedModal, ModalBody, ModalHeader, Tab, TabList, TabPanel, TabPanels, Tabs, TextInput, RadioButtonGroup, RadioButton } from "@carbon/react";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { MultiSelect } from "@carbon/react";
import { updateNotifications } from "../../../service/notifications";

export const AccessControl = () => {

  const dispatch = useAppDispatch();
  const {data: users = []} = useGetUsersQuery({});
  const {data: groups = []} = useGetGroupsQuery({});
  const {data: permissions = []} = useGetPermissionsQuery({});

  const [upsertGroup] = useUpsertGroupMutation();
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [upsertPermission] = useUpsertPermissionMutation();

  const defaultUUID = '00000000-0000-0000-0000-000000000000';
  
  const [selectedGroupName, setSelectedGroupName] = React.useState<any>('');
  const [selectedGroupID, setSelectedGroupID] = React.useState<any>(defaultUUID);
  const [selectedPermissions, setSelectedPermissions] = React.useState<any[]>([]);
  const [selectedUsers, setSelectedUsers] = React.useState<any[]>([]);

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [selectedPermissionID, setSelectedPermissionID] = React.useState<any>(defaultUUID);
  const [selectedPermissionName, setSelectedPermissionName] = React.useState<any>('');
  const [selectedPermissionAttribute, setPermissionAttribute] = React.useState<any>('rw');
  const [selectedPermissionGrant, setPermissionGrant] = React.useState<any>('');

  const [usersFilter, setUsersFilter] = React.useState('');
  const filterUsers = (args: any) => {
    setUsersFilter(args.target.value);
  };

  const [groupsFilter, setGroupsFilter] = React.useState('');
  const filterGroups = (args: any) => {
    setGroupsFilter(args.target.value);
  };

  const [permsFilter, setPermsFilter] = React.useState('');
  const filterPerms = (args: any) => {
    setPermsFilter(args.target.value);
  };


  const adminModalState = useSelector((state: RootState) => state.adminModalState);

  const setGroupsModalOpen = (open: boolean) => {
    dispatch(updateAdminModalState({openGroupsModal: open, openPermissionsModal: false}));
  };

  const setPermissionsModalOpen = (open: boolean) => {
    dispatch(updateAdminModalState({openGroupsModal: false, openPermissionsModal: open}));
  };

  const handleGroupFormSubmit = () => {
    upsertGroup({id: selectedGroupID, name: selectedGroupName, permissions: selectedPermissions, users: selectedUsers}).unwrap()
      .then(() => dispatch(updateNotifications({notifications: [{notif: 'succesful group upsert', status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error upserting group: ' + JSON.stringify(error), status: 'error'}]})));
    resetSelectedGroup();
  };

  const resetSelectedGroup = () => {
    setGroupsModalOpen(false);
    setSelectedGroupID(defaultUUID);
    setSelectedGroupName(''); 
    setSelectedPermissions([]);
    setSelectedUsers([]);
  };

  const handlePermissionFormSubmit = () => {
    const grant = selectedPermissionGrant.includes('*') ? '*' : selectedPermissionGrant + '_' + selectedPermissionAttribute;
    upsertPermission({id: selectedPermissionID, name: selectedPermissionName, grant: grant}).unwrap()
      .then(() => dispatch(updateNotifications({notifications: [{notif: 'succesful permission upsert', status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error upserting permission: ' + JSON.stringify(error), status: 'error'}]})));;
    resetSelectedPermission();
  };

  const resetSelectedPermission = () => {
    setSelectedPermissionID(defaultUUID);
    setPermissionsModalOpen(false);
    setSelectedPermissionName(''); 
    setPermissionAttribute('');
    setPermissionGrant('');
  };

  return (
    <>
      <Tabs>
          <TabList aria-label="Admin tabs">
            <Tab renderIcon={() => {return <Identification/>;}}>Users</Tab>
            <Tab renderIcon={() => {return <Events/>;}}>Groups</Tab>
            <Tab renderIcon={() => {return <GroupSecurity/>;}}>Permissions</Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <ResourceDataTable 
                rows={users.filter((perm: any) => {
                  return (
                    (perm?.name && perm.name.toLowerCase().includes(usersFilter.toLowerCase())) || 
                    (perm?.email && perm.email.toLowerCase().includes(usersFilter.toLowerCase()))
                  );
                }).map((user: any) => {
                  return {
                    id: user.id,
                    name: user.name, 
                    email: user.email, 
                    isAdmin: user.isAdmin === true ? <CheckmarkFilled color="green" /> : <Misuse color="coral"/>
                  };
                })}
                headers={[{'header': 'Name', 'key': 'name'}, {'header': 'Email', 'key': 'email'}, {'header': 'Admin', 'key': 'isAdmin'}]} 
                filterFunction={filterUsers}
                filterPlaceholder="Filter users"
                filterValue={usersFilter}
                title={'Users'}
                batchActions={[]}
              />
            </TabPanel>
            <TabPanel>
            <AdminDataTable 
                rows={groups.filter((group: any) => {
                  return (
                    (group?.name && group.name.toLowerCase().includes(groupsFilter.toLowerCase()))
                  );
                }).map((group: any) => 
                  {
                    return {
                      id: group.id,
                      name: group.name, 
                      permissions: group.permissions.map((perm: any) => {return <Tag key={perm.name} type="blue">{perm.name + ' (' + perm.appTag + ')'}</Tag> ;}),
                      users: group.users.map((user: any) => {return <Tag key={user.name} type="green">{user.name}</Tag>;}),
                      actions: <ButtonSet style={{maxWidth: '50px'}}>
                            <Button 
                              renderIcon={Edit}
                              kind="ghost" 
                              size="sm"
                              onClick={() => {
                                setSelectedGroupName(group.name); 
                                setSelectedGroupID(group.id); 
                                setSelectedPermissions(group.permissions); 
                                setSelectedUsers(group.users); 
                                setGroupsModalOpen(true);
                              }}
                            >
                              Edit
                            </Button> 
                            <Button 
                              disabled
                              renderIcon={Close}
                              kind="ghost" 
                              size="sm"
                            >
                              Delete
                            </Button>
                        </ButtonSet>
                    };
                  })
                }
                headers={[{'header': 'Name', 'key': 'name'}, {'header': 'Users', 'key': 'users'}, {'header': 'Permissions', 'key': 'permissions'}, {'header': 'Actions', 'key': 'actions'}]} 
                filterFunction={filterGroups} 
                filterPlaceholder="Filter groups"
                filterValue={groupsFilter}
                title={'Groups'}
                upsertFunction={() => setGroupsModalOpen(true)}
                upsertFunctionTitle={'Create Group'}
              /> 
            </TabPanel>
            <TabPanel>
            <AdminDataTable 
                rows={permissions.filter((perm: any) => {
                  return (
                    (perm?.name && perm.name.toLowerCase().includes(permsFilter.toLowerCase())) || 
                    (perm?.appTag && perm.appTag.toLowerCase().includes(permsFilter.toLowerCase()))
                  );
                }).map((perm: any) =>
                  {
                    return {
                      id: perm.id,
                      name: perm.name, 
                      grant: <Tag key={perm.name} type="blue">{perm.appTag}</Tag>,
                      actions: <ButtonSet style={{maxWidth: '50px'}}>
                            <Button 
                              renderIcon={Edit} 
                              kind="ghost" 
                              size="sm"
                              onClick={() => {
                                setSelectedPermissionID(perm.id); 
                                setSelectedPermissionName(perm.name);
                                setPermissionGrant(perm.appTag.replace(/_write|_read/g, ''));
                                setPermissionAttribute(perm.appTag.includes('_write') ? 'write' : 'read');
                                setPermissionsModalOpen(true);
                              }}
                            >
                              Edit
                            </Button> 
                            <Button 
                              disabled
                              renderIcon={Close} 
                              kind="ghost" 
                              size="sm"
                            >
                              Delete
                            </Button>
                        </ButtonSet>
                    };
                  })
                }
                headers={[{'header': 'Name', 'key': 'name'}, {'header': 'Grant', 'key': 'grant'}, {'header': 'Actions', 'key': 'actions'}]} 
                filterFunction={filterPerms} 
                filterPlaceholder="Filter permissions"
                filterValue={permsFilter}
                title={'Permissions'}
                upsertFunction={() => setPermissionsModalOpen(true)}
                upsertFunctionTitle={'Create Permission'}
            />
          </TabPanel>
        </TabPanels>
      </Tabs>

      <ComposedModal open={adminModalState.openGroupsModal} onClose={() => {resetSelectedGroup();}}>
        <ModalHeader label="Groups" title="Add a new group" />
        <ModalBody>
          <p style={{
          marginBottom: '1rem'
        }}>
            Groups are used to manage permissions for a set of users. Add or remove users and permissions from the group.
          </p>
          <TextInput 
            data-modal-primary-focus 
            onChange={(e: any) => {setSelectedGroupName(e.target.value);}}
            id="text-input-1" 
            labelText="Group Name" 
            placeholder="e.g. Grid" 
            style={{marginBottom: '1rem'}}
            value={selectedGroupName} 
          />
          <MultiSelect 
            id="perm-select" 
            label="" 
            titleText="Select permissions"
            items={permissions.map((perm: any) => {return {id: perm.id, name: perm.name, appTag: perm.appTag};})}
            selectedItems={selectedPermissions}
            itemToString={(item: any) => item.name + ' (' + item.appTag + ')'}
            onChange={(args: any) => { setSelectedPermissions(args.selectedItems); }}
          />
          <div style={{marginBottom: '20px'}}></div>
          <MultiSelect
            id="user-select" 
            label="" 
            titleText="Select users"
            items={users.map((user: any) => {return {id: user.id, name: user.name, email: user.email};})}
            selectedItems={selectedUsers}
            itemToString={(item: any) => item.name}
            onChange={(args: any) => { setSelectedUsers(args.selectedItems); }}
          />
          <ButtonSet style={{marginTop: '20px'}}>
            <Button kind="primary" onClick={() => handleGroupFormSubmit()}>
              Submit
            </Button>
            <Button kind="secondary" onClick={() => {resetSelectedGroup();}}>
              Cancel
            </Button>
          </ButtonSet>
        </ModalBody>
      </ComposedModal>

      <ComposedModal open={adminModalState.openPermissionsModal} onClose={() => {resetSelectedPermission();}}>
        <ModalHeader label="Permissions" title="Add a new permission" />
        <ModalBody>
          <p style={{marginBottom: '1rem'}}>
            Permissions are used to manage access to k8ds resources based on their labels.
          </p>
          <TextInput 
            data-modal-primary-focus 
            onChange={(e: any) => {setSelectedPermissionName(e.target.value);}}
            id="text-input-1" 
            labelText="Permission Name" 
            placeholder="e.g. GridBatch" 
            style={{marginBottom: '1rem'}}
            value={selectedPermissionName} 
          />
          <TextInput 
            data-modal-primary-focus 
            onChange={(e: any) => {setPermissionGrant(e.target.value);}}
            id="text-input-2" 
            labelText="Permission grant" 
            placeholder="e.g. batch-grid-c --> grants access to batch-grid-c resources" 
            style={{marginBottom: '1rem'}}
            value={selectedPermissionGrant} 
          />
         
          <RadioButtonGroup 
            valueSelected={selectedPermissionAttribute} 
            value={selectedPermissionAttribute}
            onChange={(e: any) => {setPermissionAttribute(e);} }
            legendText="Select permission access level" 
            name="attribute-selector" 
            defaultSelected="radio-1">
              <RadioButton labelText="read" value="read" id="radio-1"/>
              <RadioButton labelText="write" value="write" id="radio-2"/>
          </RadioButtonGroup>
          <ButtonSet style={{marginTop: '20px'}}>
            <Button kind="primary" onClick={() => handlePermissionFormSubmit()}>
              Submit
            </Button>
            <Button kind="secondary" onClick={() => {resetSelectedPermission();}}>
              Cancel
            </Button>
          </ButtonSet>
        </ModalBody>
      </ComposedModal>
    </>
  );
};