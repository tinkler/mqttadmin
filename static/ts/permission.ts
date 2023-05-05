// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
import { Role } from './role';
import { User } from './user';


export interface Permission {
	
	id: number;
	
	name: string;
	
	
	savePermission(): Promise<void>;
	
	deletePermission(): Promise<void>;
	
}

export interface RolePermission {
	
	id: number;
	
	role: Role;
	
	permission: Permission;
	
	
	savePermission(): Promise<void>;
	
	deletePermission(): Promise<void>;
	
}

export interface UserPermission {
	
	id: number;
	
	user: User;
	
	permission: Permission;
	
	
	savePermission(): Promise<void>;
	
	deletePermission(): Promise<void>;
	
}



export function Permission(): Permission {
	
	return {
		
		

		id: 0,
		

		name: "",
		

		

		savePermission(): Promise<void> {
			
			return postPermission(this, 'save-permission', {  });
			
		},
		

		deletePermission(): Promise<void> {
			
			return postPermission(this, 'delete-permission', {  });
			
		},
		
		
	};
	
}

// post data by restful api

function postPermission(permission: Permission, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/permission/permission/${method}`, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	return new Promise((resolve, reject) => {
		xhr.onload = () => {
			if (xhr.status === 200) {
				resolve(xhr.response);
			} else {
				reject(new Error(xhr.statusText));
			}
		};
		xhr.onerror = () => {
			reject(new Error(xhr.statusText));
		};
		xhr.send(JSON.stringify({ data: permission, args }));
	});
}

export function RolePermission(): RolePermission {
	
	return {
		
		

		id: 0,
		

		role: Role(),
		

		permission: Permission(),
		

		

		savePermission(): Promise<void> {
			
			return postRolePermission(this, 'save-permission', {  });
			
		},
		

		deletePermission(): Promise<void> {
			
			return postRolePermission(this, 'delete-permission', {  });
			
		},
		
		
	};
	
}

// post data by restful api

function postRolePermission(rolePermission: RolePermission, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/permission/role_permission/${method}`, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	return new Promise((resolve, reject) => {
		xhr.onload = () => {
			if (xhr.status === 200) {
				resolve(xhr.response);
			} else {
				reject(new Error(xhr.statusText));
			}
		};
		xhr.onerror = () => {
			reject(new Error(xhr.statusText));
		};
		xhr.send(JSON.stringify({ data: rolePermission, args }));
	});
}

export function UserPermission(): UserPermission {
	
	return {
		
		

		id: 0,
		

		user: User(),
		

		permission: Permission(),
		

		

		savePermission(): Promise<void> {
			
			return postUserPermission(this, 'save-permission', {  });
			
		},
		

		deletePermission(): Promise<void> {
			
			return postUserPermission(this, 'delete-permission', {  });
			
		},
		
		
	};
	
}

// post data by restful api

function postUserPermission(userPermission: UserPermission, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/permission/user_permission/${method}`, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	return new Promise((resolve, reject) => {
		xhr.onload = () => {
			if (xhr.status === 200) {
				resolve(xhr.response);
			} else {
				reject(new Error(xhr.statusText));
			}
		};
		xhr.onerror = () => {
			reject(new Error(xhr.statusText));
		};
		xhr.send(JSON.stringify({ data: userPermission, args }));
	});
}

