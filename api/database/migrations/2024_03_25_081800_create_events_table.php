<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('events', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->string('state');
            $table->string('on_home');
            $table->string('title');
            $table->string('geolink');
            $table->string('type');
            $table->string('location')->nullable();
            $table->dateTime('start_date');
            $table->dateTime('end_date');
            $table->string('imageUrl');
            $table->string('ticket_purchase_url')->nullable();

            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('events');
    }
};
