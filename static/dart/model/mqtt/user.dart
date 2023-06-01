// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
import '../../http.dart';
import './const.dart';
import './role.dart' as $role show Role;


class Auth {
	
	String id = "";
	
	String deviceToken = "";
	
	String username = "";
	
	String password = "";
	
	String token = "";
	
	
	/// Signin sign in with username and password
	/// Require: Username, Password
	/// Optional: DeviceToken
	Future<Auth?> signin(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/auth/signin', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = Auth.fromJson(response.data['data']['data']);
			assign(respModel);
			if (response.data['data']['resp'] != null) {
				return Auth.fromJson(response.data['data']['resp']);
			} else {
				return null;
			}
			
		}
		return null;
		
	}
	
	/// QuickSignin quick signin without password
	/// Require: DeviceToken, Token
	Future<void> quickSignin(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/auth/quick-signin', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = Auth.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	Future<Auth?> signup(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/auth/signup', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = Auth.fromJson(response.data['data']['data']);
			assign(respModel);
			if (response.data['data']['resp'] != null) {
				return Auth.fromJson(response.data['data']['resp']);
			} else {
				return null;
			}
			
		}
		return null;
		
	}
	
	Auth();

	assign(Auth other) {
		
		id = other.id;
		
		deviceToken = other.deviceToken;
		
		username = other.username;
		
		password = other.password;
		
		token = other.token;
		
	}

	Map<String, dynamic> toJson() {
		return {
			
			"id": id,
			
			"device_token": deviceToken,
			
			"username": username,
			
			"password": password,
			
			"token": token,
			
		};
	}
	Auth.fromJson(Map<String, dynamic> json) {
		
		id = json["id"];
		
		deviceToken = json["device_token"];
		
		username = json["username"];
		
		password = json["password"];
		
		token = json["token"];
		
	}
}

class UserProfile {
	
	String phoneNo = "";
	
	Future<void> save(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/user_profile/save', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = UserProfile.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	UserProfile();

	assign(UserProfile other) {
		
		phoneNo = other.phoneNo;
		
	}

	Map<String, dynamic> toJson() {
		return {
			
			"phone_no": phoneNo,
			
		};
	}
	UserProfile.fromJson(Map<String, dynamic> json) {
		
		phoneNo = json["phone_no"];
		
	}
}

class User {
	
	String id = "";
	
	String username = "";
	
	String email = "";
	
	List<UserProfile> profiles = [];
	
	List<$role.Role> roles = [];
	
	
	/// Save saves the user to the database
	Future<void> save(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/user/save', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = User.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	/// AddRole adds a role to the user
	Future<void> addRole(
		$role.Role? role,
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/user/add-role', data: {
			"data": toJson(),
			"args": { "role": role, }
		});
		if (response.data['code'] == 0) {
			var respModel = User.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	Future<void> removeRole(
		$role.Role? role,
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/user/remove-role', data: {
			"data": toJson(),
			"args": { "role": role, }
		});
		if (response.data['code'] == 0) {
			var respModel = User.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	/// Get gets the user from the database
	/// Only admin can get other user's information
	Future<void> get(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/user/get', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = User.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	/// GetRoles gets the roles of the user
	Future<void> getRoles(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/user/get-roles', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = User.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	User();

	assign(User other) {
		
		id = other.id;
		
		username = other.username;
		
		email = other.email;
		
		profiles = other.profiles;
		
		roles = other.roles;
		
	}

	Map<String, dynamic> toJson() {
		return {
			
			"id": id,
			
			"username": username,
			
			"email": email,
			
			"profiles": profiles.map((e) => e.toJson()).toList(),
			
			"roles": roles.map((e) => e.toJson()).toList(),
			
		};
	}
	User.fromJson(Map<String, dynamic> json) {
		
		id = json["id"];
		
		username = json["username"];
		
		email = json["email"];
		
		profiles = json["profiles"] == null ? [] : (json["profiles"] as List<dynamic>).map((e) => UserProfile.fromJson(e)).toList();
		
		roles = json["roles"] == null ? [] : (json["roles"] as List<dynamic>).map((e) => $role.Role.fromJson(e)).toList();
		
	}
}

class UserRole {
	
	String id = "";
	
	User? user;
	
	$role.Role? role;
	
	Future<void> save(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/user/user_role/save', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = UserRole.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	UserRole();

	assign(UserRole other) {
		
		id = other.id;
		
		user = other.user;
		
		role = other.role;
		
	}

	Map<String, dynamic> toJson() {
		return {
			
			"id": id,
			
			"user": user != null ? user!.toJson() : null,
			
			"role": role != null ? role!.toJson() : null,
			
		};
	}
	UserRole.fromJson(Map<String, dynamic> json) {
		
		id = json["id"];
		
		user = json["user"] == null ? User() : User.fromJson(json["user"]);
		
		role = json["role"] == null ? $role.Role() : $role.Role.fromJson(json["role"]);
		
	}
}


