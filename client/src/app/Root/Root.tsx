import React, { useEffect } from 'react';
import { Shell } from './Shell';
import { ContentArea } from './ContentArea';
import { Route, Routes } from 'react-router-dom';
import { LoginFlow, ProtectedRoutes } from '../components/ProtectedRoutes/ProtectedRoutes';
import { useUserInfoQuery } from '../../service/khub';
import { Loading } from '@carbon/react';

import { updateNotifications } from '../../service/notifications';
import { useAppDispatch } from '../store';
import { updateUserIsAdmin } from '../../service/user-admin-state';

export type RootProps = {
  setExpanded: (expanded: boolean) => void;
  expand: boolean;
};

export function Root() {
  const dispatch = useAppDispatch();
  const {data: userInfo} = useUserInfoQuery({});
  const userPreferredTheme = userInfo?.darkMode ? 'dark' : 'light';
  

  useEffect(() => {
    if (userInfo !== undefined && (userInfo.isAdmin || 
      (userInfo.groups && userInfo.groups.some((elem: any) => { return elem.name === 'Admin'; })))) {
      dispatch(updateUserIsAdmin({isAdmin: true}));
    } else {
      dispatch(updateUserIsAdmin({isAdmin: false}));
    }
  });

  return (
    <div style={{height: '100%'}} className='App'>
      <Routes>
        <Route element={<ProtectedRoutes/>}>
          {userInfo && 
            <Route path="/*" element={
              <div>
                <Shell userPreferredTheme={userPreferredTheme} userInfo={userInfo}/>
                <ContentArea />
              </div> 
            }
            />
          }
          {!userInfo && 
            <Route path="/*" element={<Loading withOverlay={true} description='no user auth session'/>}/>
          }
          
        </Route>
        <Route path="/sign-in" element={<LoginFlow />} />
      </Routes>
    </div>
  );
}
