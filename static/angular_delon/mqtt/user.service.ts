// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
import { Role } from './role.model';

import { Injectable } from '@angular/core';
import { UserRole, Auth, User, UserProfile,  } from './user.model';
import { RoleService  } from './role.service';

import { _HttpClient } from '@delon/theme';

@Injectable({
	providedIn: 'root'
})
export class UserService {
  
	constructor(
		private http: _HttpClient,
		private roleSrv: RoleService,
		
		) { }

	
	newUserRole(): UserRole {
		return new UserRole(this.http);
	}
	newAuth(): Auth {
		return new Auth(this.http);
	}
	/**
	* User is the user model
	*/
	newUser(): User {
		return new User(this.http);
	}
	newUserProfile(): UserProfile {
		return new UserProfile(this.http);
	}
}
