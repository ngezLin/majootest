<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;

return new class extends Migration
{
    public function up(): void
    {
        // Create Users table matching Merchants.user_id foreign relation
        Schema::create('Users', function (Blueprint $table) {
            $table->id();
            $table->string('email', 100)->unique();
            $table->string('password', 255);
            $table->timestamps();
        });

        // Seed initial users for Merchant 1 (user_id 1) and Merchant 2 (user_id 2)
        DB::table('Users')->insert([
            [
                'id'         => 1,
                'email'      => 'merchant1@mail.com',
                'password'   => Hash::make('password123'),
                'created_at' => now(),
                'updated_at' => now(),
            ],
            [
                'id'         => 2,
                'email'      => 'merchant2@mail.com',
                'password'   => Hash::make('password123'),
                'created_at' => now(),
                'updated_at' => now(),
            ],
        ]);
    }

    public function down(): void
    {
        Schema::dropIfExists('Users');
    }
};
