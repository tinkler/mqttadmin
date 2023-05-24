// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
import { Role } from './role';


export interface Auth {
	
	/**
	* UUID
	*/
	id: string;
	
	/**
	* UUID
	*/
	deviceToken: string;
	
	username: string;
	
	password: string;
	
	token: string;
	
	
	signin(): Promise<Auth> ;
	
	/**
	* QuickSignin quick signin with password
	*/
	quickSignin(): Promise<void>;
	
	signup(): Promise<Auth> ;
	
}

/**
* User is the user model
*/
export interface User {
	
	/**
	* ID is the primary key
	*/
	id: string;
	
	username: string;
	
	email: string;
	
	profiles: UserProfile[];
	
	
	/**
	* Save saves the user to the database
	*/
	save(): Promise<void>;
	
	/**
	* AddRole adds a role to the user
	*/
	addRole(role: Role, ): Promise<void>;
	
}

export interface UserProfile {
	
	phoneNo: string;
	
	
	save(): Promise<void>;
	
}

export interface UserRole {
	
	id: number;
	
	user: User;
	
	role: Role;
	
	
	save(): Promise<void>;
	
}



export function Auth(): Auth {
	
	return {
		
		id: "",
		
		deviceToken: "",
		
		username: "",
		
		password: "",
		
		token: "",
		
		
		signin(): Promise<Auth>  {
			
			return postAuth(this, 'signin', {  }).then((res: { data: any }) => res.data as Auth);
			
		},
		
		quickSignin(): Promise<void> {
			
			return postAuth(this, 'quick-signin', {  });
			
		},
		
		signup(): Promise<Auth>  {
			
			return postAuth(this, 'signup', {  }).then((res: { data: any }) => res.data as Auth);
			
		},
		
		
	};
	
}

// post data by restful api

function postAuth(auth: Auth, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/user/auth/${method}`, true);
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
		xhr.send(JSON.stringify({ data: auth, args }));
	});
}

/**
* User is the user model
*/
export function User(): User {
	
	return {
		
		id: "",
		
		username: "",
		
		email: "",
		
		profiles: ,
		
		
		save(): Promise<void> {
			
			return postUser(this, 'save', {  });
			
		},
		
		addRole(role: Role, ): Promise<void> {
			
			return postUser(this, 'add-role', { role,  });
			
		},
		
		
	};
	
}

// post data by restful api

function postUser(user: User, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/user/user/${method}`, true);
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
		xhr.send(JSON.stringify({ data: user, args }));
	});
}

export function UserProfile(): UserProfile {
	
	return {
		
		phoneNo: "",
		
		
		save(): Promise<void> {
			
			return postUserProfile(this, 'save', {  });
			
		},
		
		
	};
	
}

// post data by restful api

function postUserProfile(userProfile: UserProfile, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/user/user_profile/${method}`, true);
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
		xhr.send(JSON.stringify({ data: userProfile, args }));
	});
}

export function UserRole(): UserRole {
	
	return {
		
		id: 0,
		
		user: User(),
		
		role: Role(),
		
		
		save(): Promise<void> {
			
			return postUserRole(this, 'save', {  });
			
		},
		
		
	};
	
}

// post data by restful api

function postUserRole(userRole: UserRole, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/user/user_role/${method}`, true);
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
		xhr.send(JSON.stringify({ data: userRole, args }));
	});
}

